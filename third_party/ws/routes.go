package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager quản lý tất cả các kết nối WebSocket và phòng
type ConnectionManager struct {
	clients    map[*websocket.Conn]bool
	roomUsers  map[int]map[*websocket.Conn]int // roomID -> connections -> userID
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

// NewConnectionManager tạo một ConnectionManager mới
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients:    make(map[*websocket.Conn]bool),
		roomUsers:  make(map[int]map[*websocket.Conn]int),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run khởi chạy các goroutine quản lý kết nối
func (manager *ConnectionManager) Run() {
	for {
		select {
		case conn := <-manager.register:
			manager.mu.Lock()
			manager.clients[conn] = true
			manager.mu.Unlock()
			fmt.Println("Client đã kết nối!")

		case conn := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[conn]; ok {
				// Xóa khỏi danh sách kết nối
				delete(manager.clients, conn)
				manager.removeFromRooms(conn)
				conn.Close()
				fmt.Println("Client đã ngắt kết nối!")
			}
			manager.mu.Unlock()

		case message := <-manager.broadcast:
			manager.mu.Lock()
			// Gửi tin nhắn đến tất cả client
			for conn := range manager.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					fmt.Println("Lỗi khi gửi tin nhắn:", err)
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mu.Unlock()
		}
	}
}

// AddToRoom thêm người dùng vào một phòng
func (manager *ConnectionManager) AddToRoom(roomID, userID int, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if _, exists := manager.roomUsers[roomID]; !exists {
		manager.roomUsers[roomID] = make(map[*websocket.Conn]int)
	}
	manager.roomUsers[roomID][conn] = userID
}

// RemoveFromRoom xóa người dùng khỏi một phòng
func (manager *ConnectionManager) RemoveFromRoom(roomID int, conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if users, exists := manager.roomUsers[roomID]; exists {
		delete(users, conn)
		if len(users) == 0 {
			delete(manager.roomUsers, roomID)
		}
	}
}

// removeFromRooms xóa kết nối khỏi tất cả các phòng
func (manager *ConnectionManager) removeFromRooms(conn *websocket.Conn) {
	for roomID, users := range manager.roomUsers {
		if _, exists := users[conn]; exists {
			delete(users, conn)
			if len(users) == 0 {
				delete(manager.roomUsers, roomID)
			}
		}
	}
}

// HandleConnections xử lý các kết nối WebSocket
func HandleConnections(manager *ConnectionManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Lỗi khi nâng cấp:", err)
		return
	}
	manager.register <- conn

	defer func() {
		manager.unregister <- conn
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			manager.unregister <- conn
			break
		}

		var msgData map[string]interface{}
		if err := json.Unmarshal(message, &msgData); err != nil {
			log.Println("Error unmarshal message", err)
			continue
		}

		// Xử lý thông điệp theo logic
		manager.handleMessage(msgData, conn)
	}
}

// handleMessage xử lý thông điệp từ client
func (manager *ConnectionManager) handleMessage(msgData map[string]interface{}, conn *websocket.Conn) {
	roomID := int(msgData["room_id"].(float64))
	userID := int(msgData["user_id"].(float64))
	action := msgData["action"].(string)
	msgContent := msgData["message"].(string)

	switch action {
	case "join":
		manager.AddToRoom(roomID, userID, conn)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("User %d joined room %d", userID, roomID)))

	case "leave":
		manager.RemoveFromRoom(roomID, conn)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("User %d left room %d", userID, roomID)))

	case "send_message":
		manager.mu.Lock()
		if users, exists := manager.roomUsers[roomID]; exists {
			for userConn := range users {
				userConn.WriteMessage(websocket.TextMessage, []byte(msgContent))
			}
		}
		manager.mu.Unlock()

	default:
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid action"))
	}
}

// Cấu hình upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
