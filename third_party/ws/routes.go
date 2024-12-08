package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép tất cả origin (xem xét nếu cần bảo mật)
	},
}

type ConnectionManager struct {
	mu          sync.Mutex
	connections map[int]map[int]*websocket.Conn
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[int]map[int]*websocket.Conn),
	}
}

// Thay thế kết nối cũ nếu tồn tại
func (cm *ConnectionManager) ReplaceConnectionIfExists(roomId int, userId int, conn *websocket.Conn) {
	if existingConn := cm.connections[roomId][userId]; existingConn != nil {
		existingConn.Close()
	}
	cm.AddConnection(roomId, userId, conn)
}

func (cm *ConnectionManager) AddConnection(roomId int, userId int, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.connections[roomId]; !ok {
		cm.connections[roomId] = make(map[int]*websocket.Conn)
	}
	cm.connections[roomId][userId] = conn
}

func (cm *ConnectionManager) RemoveConnection(roomId int, userId int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if users, ok := cm.connections[roomId]; ok {
		delete(users, userId)
		if len(users) == 0 {
			delete(cm.connections, roomId)
		}
	}
}

func (cm *ConnectionManager) BroadcastToRoom(roomId int, message []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if users, ok := cm.connections[roomId]; ok {
		for userId, conn := range users {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				global.Logger.Sugar().Error("Error sending message to user %d in room %d: %v", userId, roomId, err)
			}
		}
	}

}

func HandleConnection(w http.ResponseWriter, r *http.Request, cm *ConnectionManager) {

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "havnt gin.Context", http.StatusUnauthorized)
		return
	}
	userId := auth.GetUserIdFromToken(w, tokenString)
	global.Logger.Sugar().Infof("User %s connected", userId)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Nếu là yêu cầu OPTIONS, trả về HTTP 200 OK để cho phép CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		global.Logger.Sugar().Error("Failed to upgrade connection", err)
		return
	}
	defer conn.Close()

	const (
		pongWait   = 60 * time.Second
		pingPeriod = (pongWait * 9) / 10
	)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		global.Logger.Sugar().Info("Pong received from client")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				global.Logger.Sugar().Error("Failed to send Ping", err)
				conn.Close()
				return
			}
		}
	}()

	// Lặp vô hạn để nhận tin nhắn từ client
	for {
		global.Logger.Sugar().Info("Reading message from WebSocket...")
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Sugar().Errorf("WebSocket connection closed unexpectedly (userId: %d): %v", userId, err)
			} else {
				global.Logger.Sugar().Errorf("Error reading WebSocket message (userId: %d): %v", userId, err)
			}
			return
		}
		global.Logger.Sugar().Infof("Received message: %s (userId: %s)", string(message), userId)

		var msgData map[string]interface{}
		if err := json.Unmarshal(message, &msgData); err != nil {
			global.Logger.Sugar().Error("Error Unmarshal message", err)
			continue
		}

		// xac thuc jwt trong tung tin nhan
		msgToken, ok := msgData["token"].(string)
		if !ok || msgToken == "" {
			global.Logger.Sugar().Error("Missing or invalid token in message")
			err = conn.WriteMessage(messageType, []byte("Unauthorized: missing token"))
			if err != nil {
				break
			}
			continue
		}

		userId := auth.GetUserIdFromToken(w, msgToken)
		// lay thong tin
		// Lấy và chuyển đổi giá trị room_id
		roomIdFloat, ok := msgData["room_id"].(float64)
		if !ok {
			log.Println("room_id không hợp lệ hoặc không tồn tại")
			return
		}

		// Chuyển đổi từ float64 sang int
		roomId := int(roomIdFloat)

		action := msgData["action"].(string)
		msgContent := msgData["message"].(string)
		// Thực hiện hành động tùy thuộc vào yêu cầu
		switch action {
		case "join":
			cm.ReplaceConnectionIfExists(roomId, userId, conn)
			cm.AddConnection(roomId, userId, conn)
			// Người dùng tham gia nhóm
			result := handleJoinGroup(userId, roomId)
			err = conn.WriteMessage(messageType, []byte(result))
		case "leave":
			// Người dùng rời nhóm
			cm.RemoveConnection(userId, roomId)
			result := handleLeaveGroup(userId, roomId)

			err = conn.WriteMessage(messageType, []byte(result))
		case "send_message":
			// Người dùng gửi tin nhắn
			handleSendMessage(userId, roomId, msgContent)
			err = conn.WriteMessage(messageType, []byte(msgContent))
		default:
			// Nếu không phải hành động hợp lệ
			err = conn.WriteMessage(messageType, []byte("Invalid action"))
		}

		// Lắng nghe Redis stream trong một goroutine riêng
		go listenRedisStream(roomId, cm)
		// Kiểm tra lỗi khi gửi phản hồi tới client
		if err != nil {
			global.Logger.Sugar().Error("Error sending response:", err)
			break
		}
	}
}

func listenRedisStream(roomId int, cm *ConnectionManager) {

	streamKey := fmt.Sprintf("room:%d:messages", roomId)

	// Lắng nghe stream mà không sử dụng consumer group
	for {
		messages, err := global.Rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{streamKey, "$"},
			Block:   0,
			Count:   1,
		}).Result()

		if err != nil {
			log.Printf("Error reading message from Redis Stream for room %d: %v", roomId, err)
			break
		}

		for _, message := range messages {
			for _, xMessage := range message.Messages {
				// Gửi tin nhắn tới WebSocket
				msg, ok := xMessage.Values["message"].(string)
				if !ok {
					log.Printf("Invalid message format for room %d: %v", roomId, xMessage.Values)
					continue
				}
				cm.BroadcastToRoom(roomId, []byte(msg))

				// Không cần xác nhận như khi sử dụng consumer group
			}
		}
	}

}

func handleJoinGroup(userId int, roomId int) (result string) {
	// Cập nhật MySQL: Thêm người dùng vào bảng Memberships
	global.Logger.Sugar().Info(userId)
	queries := database.New(global.Mdbc)
	params := database.AddMemberToRoomChatParams{
		UserID: uint64(userId),
		RoomID: int64(roomId),
	}
	err := queries.AddMemberToRoomChat(context.Background(), params)
	if err != nil {
		log.Println("Error adding user to room in MySQL:", err)
		return
	}

	// Cập nhật Redis: Thêm người dùng vào danh sách thành viên trong phòng
	err = global.Rdb.SAdd(context.Background(), "room:"+strconv.Itoa(roomId)+":members", userId).Err()
	if err != nil {
		log.Println("Error adding user to room in Redis:", err)
		return
	}
	err = global.Rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "room:" + strconv.Itoa(roomId) + ":messages",
		Values: map[string]interface{}{"user_id": userId, "message": "user join room"},
	}).Err()
	if err != nil {
		log.Println("Error adding message to Redis stream:", err)
		return
	}
	result = fmt.Sprintf("UserId %d join group", userId)
	return result

}

func handleLeaveGroup(userId int, roomId int) (result string) {
	queries := database.New(global.Mdbc)
	// Cập nhật MySQL: Xóa người dùng khỏi bảng Memberships
	params := database.DeleteMemberFromRoomChatParams{
		UserID: uint64(userId),
		RoomID: int64(roomId),
	}
	err := queries.DeleteMemberFromRoomChat(context.Background(), params)
	if err != nil {
		log.Println("Error removing user from room in MySQL:", err)
		return
	}

	// Cập nhật Redis: Xóa người dùng khỏi danh sách thành viên trong phòng
	err = global.Rdb.SRem(context.Background(), "room:"+strconv.Itoa(roomId)+":members", userId).Err()
	if err != nil {
		log.Println("Error removing user from room in Redis:", err)
		return
	}

	err = global.Rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "room:" + strconv.Itoa(roomId) + ":messages",
		Values: map[string]interface{}{"user_id": userId, "message": "user leave room"},
	}).Err()
	if err != nil {
		log.Println("Error adding message to Redis stream:", err)
		return
	}
	result = fmt.Sprintf("UserId %d leave group", userId)
	return result
}

func handleSendMessage(userId int, roomId int, message string) {
	// Gửi tin nhắn vào Redis Stream
	err := global.Rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "room:" + strconv.Itoa(roomId) + ":messages",
		Values: map[string]interface{}{"user_id": userId, "message": message},
	}).Err()
	if err != nil {
		log.Println("Error adding message to Redis stream:", err)
		return
	}

}
