package websocket

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service/impl"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserInfo *model.UserInfo
	Conn     *websocket.Conn
	Auth     bool
}

type ConnectionManager struct {
	clients    map[*websocket.Conn]*Client         // Quản lý tất cả các kết nối
	roomUsers  map[int]map[*websocket.Conn]*Client // roomID -> danh sách Client
	broadcast  chan []byte                         // Kênh phát tin nhắn
	register   chan *Client                        // Kênh đăng ký Client
	unregister chan *websocket.Conn                // Kênh hủy đăng ký Client
	mu         sync.Mutex                          // Đồng bộ hóa
}

// Hàm khởi tạo ConnectionManager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients:    make(map[*websocket.Conn]*Client),
		roomUsers:  make(map[int]map[*websocket.Conn]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *websocket.Conn),
	}
}
func (cm *ConnectionManager) Run() {
	go func() {
		for {
			select {
			case client := <-cm.register:

				cm.clients[client.Conn] = client

			case conn := <-cm.unregister:
				cm.RemoveClient(conn)

			case message := <-cm.broadcast:
				// Phát tin nhắn tới tất cả các phòng
				for roomID := range cm.roomUsers {
					go cm.BroadcastToRoom(roomID, message)
				}
			}
		}
	}()
}

func (cm *ConnectionManager) AddClient(client *Client, roomID int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	// Thêm client vào phòng
	if _, exists := cm.roomUsers[roomID]; !exists {
		cm.roomUsers[roomID] = make(map[*websocket.Conn]*Client)
	}
	cm.roomUsers[roomID][client.Conn] = client
}

// RemoveFromRoom xóa người dùng khỏi một phòng
func (cm *ConnectionManager) RemoveClient(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Xóa client khỏi danh sách tổng client
	if client, exists := cm.clients[conn]; exists {
		delete(cm.clients, conn)
		// xoa client trong ds tong room
		for roomId, users := range cm.roomUsers {
			if _, exists := users[conn]; exists {
				delete(users, conn)
				if len(users) == 0 {
					delete(cm.roomUsers, roomId)
				}
				break
			}
		}
		client.Conn.Close()
	}
}

func (cm *ConnectionManager) BroadcastToRoom(roomID int, message []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if clients, exists := cm.roomUsers[roomID]; exists {
		for _, client := range clients {
			// Ghi lại thông điệp trước khi gửi
			global.Logger.Sugar().Infof("Sending message to room %d: %s", roomID, string(message))

			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				client.Conn.Close()
				cm.RemoveClient(client.Conn)
			}
		}
	} else {
		global.Logger.Sugar().Infof("No clients in room %d to send message", roomID) // Thêm log nếu không có client
	}
}

func (cm *ConnectionManager) checkAuth(client *Client) bool {
	if !client.Auth {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Authentication required"))
		return false
	}
	return true
}

// HandleConnections xử lý các kết nối WebSocket
func HandleConnections(ctx *gin.Context, cm *ConnectionManager) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("Lỗi khi nâng cấp:", err)
		return
	}
	defer conn.Close()
	defer func() {
		cm.unregister <- conn
	}()

	client := &Client{Conn: conn}
	cm.register <- client

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

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Sugar().Warnf("Unexpected WebSocket closure: %v", err)
			} else {
				global.Logger.Sugar().Infof("Normal WebSocket closure or error: %v", err)
			}

			cm.unregister <- conn
			return
		}

		// Xử lý thông điệp theo logic
		cm.handleMessage(message, client, ctx)
	}
}

// handleMessage xử lý thông điệp từ client
func (cm *ConnectionManager) handleMessage(message []byte, client *Client, ctx *gin.Context) {
	var msgData map[string]interface{}
	if err := json.Unmarshal(message, &msgData); err != nil {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid message formated"))
		return
	}
	action := msgData["action"].(string)

	switch action {
	case "auth":
		tokenString := msgData["token"].(string)

		userInfo := auth.GetUserInfoFromToken(tokenString) // Xử lý token
		if userInfo.UserID == 0 {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Authentication fail"))
			return
		}

		global.Logger.Sugar().Infof("User %d authenticated", userInfo.UserID)

		client.UserInfo = &userInfo
		client.Auth = true
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Authentication successfully"))

	case "join":
		if !cm.checkAuth(client) {
			return
		}
		roomID := int(msgData["room_id"].(float64))
		cm.AddClient(client, roomID)

	case "leave":
		if !cm.checkAuth(client) {
			return
		}
		roomID := int(msgData["room_id"].(float64))
		cm.RemoveClient(client.Conn)
		client.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("User %s left room %d", client.UserInfo.UserAccount, roomID)))

	case "send_message":
		if !cm.checkAuth(client) {
			return
		}
		roomID := int(msgData["room_id"].(float64))
		typeMessage := msgData["message_type"].(string)
		content := msgData["message"].(string)

		// Tạo đối tượng ModelChat
		chatMessage := model.ModelChat{
			UserNickname:   client.UserInfo.UserNickname.String,
			MessageContext: sql.NullString{String: content, Valid: true},
			MessageType:    model.MessagesMessageType(typeMessage),
			RoomId:         sql.NullInt32{Int32: int32(roomID), Valid: true},
			CreatedAt:      sql.NullTime{Time: time.Now(), Valid: true}, // Thêm thời gian
		}

		// Ghi vào cơ sở dữ liệu
		impl.NewsChat(database.New(global.MdbcHaproxy)).SetChatHistory(ctx, &chatMessage)

		// Tạo thông điệp JSON để phát đến các client
		messageToSend := model.ModelChat{
			UserNickname:   chatMessage.UserNickname,
			MessageContext: chatMessage.MessageContext,
			MessageType:    model.MessagesMessageType(chatMessage.MessageType),
			IsPinned:       chatMessage.IsPinned,
			CreatedAt:      chatMessage.CreatedAt,
		}
		// Chuyển đổi thông điệp thành JSON
		global.Logger.Sugar().Info(messageToSend)
		jsonMessage, err := json.Marshal(messageToSend)
		if err != nil {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Error formatting message"))
			return
		}

		// Phát tin nhắn đến tất cả client trong phòng
		cm.BroadcastToRoom(roomID, jsonMessage)

	default:
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid action"))
	}
}

// Cấu hình upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
