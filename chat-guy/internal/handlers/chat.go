package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
	"chat-guy/internal/models"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	// Map to track clients in each room
	rooms = make(map[int64]*Room)
	// Map to track private message clients
	privateClients = make(map[int64]*Client)
	roomsMutex = sync.RWMutex{}
	privateClientsMutex = sync.RWMutex{}
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	room     *Room
	userID   int64
	username string
}

type Room struct {
	id       int64
	clients  map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
	mutex    sync.Mutex
}

func newRoom(id int64) *Room {
	return &Room{
		id:         id,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.mutex.Lock()
			r.clients[client] = true
			r.mutex.Unlock()
		case client := <-r.unregister:
			r.mutex.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			r.mutex.Unlock()
		case message := <-r.broadcast:
			r.mutex.Lock()
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
			r.mutex.Unlock()
		}
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	roomID := int64(1) // Get from query params in production

	roomsMutex.Lock()
	room, ok := rooms[roomID]
	if !ok {
		room = newRoom(roomID)
		rooms[roomID] = room
		go room.run()
	}
	roomsMutex.Unlock()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		room:     room,
		userID:   claims.UserID,
		username: claims.Username,
	}

	// Register client for private messages
	privateClientsMutex.Lock()
	privateClients[claims.UserID] = client
	privateClientsMutex.Unlock()

	room.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
		privateClientsMutex.Lock()
		delete(privateClients, c.userID)
		privateClientsMutex.Unlock()
	}()

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message := models.Message{
			RoomID:    c.room.id,
			UserID:    c.userID,
			Username:  c.username,
			Content:   string(rawMessage),
			Type:      "text",
			CreatedAt: time.Now(),
		}

		// Save message to database
		db := database.GetDB()
		err = db.QueryRow(`
			INSERT INTO messages (room_id, user_id, username, content, type, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`,
			message.RoomID, message.UserID, message.Username,
			message.Content, message.Type, message.CreatedAt).Scan(&message.ID)

		if err != nil {
			log.Printf("error saving message: %v", err)
			continue
		}

		messageJSON, _ := json.Marshal(message)
		c.room.broadcast <- messageJSON
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	roomID := r.URL.Query().Get("id")

	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	// Get the client's WebSocket connection
	privateClientsMutex.RLock()
	client, ok := privateClients[claims.UserID]
	privateClientsMutex.RUnlock()

	if !ok {
		http.Error(w, "WebSocket connection not found", http.StatusBadRequest)
		return
	}

	// Remove from old room if any
	if client.room != nil {
		roomsMutex.Lock()
		if _, ok := rooms[client.room.id]; ok {
			delete(rooms[client.room.id].clients, client)
			if len(rooms[client.room.id].clients) == 0 {
				delete(rooms, client.room.id)
			}
		}
		roomsMutex.Unlock()
	}

	// Join new room
	roomsMutex.Lock()
	room, ok := rooms[int64(roomID)]
	if !ok {
		room = newRoom(int64(roomID))
		rooms[int64(roomID)] = room
		go room.run()
	}
	roomsMutex.Unlock()

	room.register <- client
	client.room = room

	w.WriteHeader(http.StatusOK)
}
