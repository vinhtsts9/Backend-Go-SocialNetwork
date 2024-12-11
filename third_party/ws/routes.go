package websocket

import (
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int
	Conn   *websocket.Conn
	Auth   bool
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
				cm.AddClient(client, client.UserID) // Giả định UserID là roomID

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

	// Thêm client vào danh sách tổng
	cm.clients[client.Conn] = client

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
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				client.Conn.Close()
				cm.RemoveClient(client.Conn)
			}
		}
	}
}

// HandleConnections xử lý các kết nối WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request, cm *ConnectionManager) {
	log.Println(1)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Lỗi khi nâng cấp:", err)
		return
	}
	defer conn.Close()
	log.Println(2)
	defer func() {
		cm.unregister <- conn
	}()

	client := &Client{Conn: conn}
	cm.register <- client
	log.Println(3)
	// const (
	// 	pongWait   = 60 * time.Second
	// 	pingPeriod = (pongWait * 9) / 10
	// )

	// conn.SetReadDeadline(time.Now().Add(pongWait))
	// conn.SetPongHandler(func(string) error {
	// 	global.Logger.Sugar().Info("Pong received from client")
	// 	conn.SetReadDeadline(time.Now().Add(pongWait))
	// 	return nil
	// })

	// ticker := time.NewTicker(pingPeriod)
	// defer ticker.Stop()

	// go func() {
	// 	for range ticker.C {
	// 		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
	// 			global.Logger.Sugar().Error("Failed to send Ping", err)
	// 			conn.Close()
	// 			return
	// 		}
	// 	}
	// }()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				global.Logger.Sugar().Errorf("WebSocket connection closed unexpectedly")
			} else {
				global.Logger.Sugar().Errorf("Error reading WebSocket message ")
			}
			cm.unregister <- conn
			return
		}
		log.Println(4)

		fmt.Println("Server received message:", string(message))

		// Xử lý thông điệp theo logic
		cm.handleMessage(w, message, client)
	}
}

// handleMessage xử lý thông điệp từ client
func (cm *ConnectionManager) handleMessage(w http.ResponseWriter, message []byte, client *Client) {
	var msgData map[string]interface{}
	if err := json.Unmarshal(message, &msgData); err != nil {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Invalid message formated"))
		return
	}
	global.Logger.Sugar().Info(msgData)
	action := msgData["action"].(string)

	switch action {
	case "auth":
		tokenString := msgData["token"].(string)

		userId := auth.GetUserIdFromToken(w, tokenString) // Xử lý token
		if userId == 0 {
			client.Conn.WriteMessage(websocket.TextMessage, []byte("Authentication failed"))
			return
		}

		global.Logger.Sugar().Infof("User %d authenticated", userId)

		client.UserID = userId
		client.Auth = true
		client.Conn.WriteMessage(websocket.TextMessage, []byte("Authentication successfully"))

	case "join":
		roomID := int(msgData["room_id"].(float64))
		cm.AddClient(client, roomID)
		client.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Joined room %d", roomID)))

	case "leave":
		roomID := int(msgData["room_id"].(float64))
		cm.RemoveClient(client.Conn)
		client.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("User %d left room %d", client.UserID, roomID)))

	case "send_message":

		roomID := int(msgData["room_id"].(float64))
		content := msgData["message"].(string)
		cm.BroadcastToRoom(roomID, []byte(fmt.Sprintf("User %d: %s", client.UserID, content)))

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
