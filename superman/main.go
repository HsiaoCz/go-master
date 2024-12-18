package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// Message holds the message details
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

// ChatRoom maintains the state of the chat room
type ChatRoom struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	mu        sync.Mutex
}

// NewChatRoom creates a new ChatRoom instance
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
}

// AddClient adds a new client to the chat room
func (c *ChatRoom) AddClient(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[conn] = true
}

// RemoveClient removes a client from the chat room
func (c *ChatRoom) RemoveClient(conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, conn)
	conn.Close()
}

// BroadcastMessage sends a message to all connected clients
func (c *ChatRoom) BroadcastMessage(msg Message) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for client := range c.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			c.RemoveClient(client)
		}
	}
}

// Run starts the chat room message broadcasting loop
func (c *ChatRoom) Run() {
	for {
		msg := <-c.broadcast
		c.BroadcastMessage(msg)
	}
}

func main() {
	chatRoom := NewChatRoom()
	go chatRoom.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}
		defer chatRoom.RemoveClient(conn)

		// Add the client to the chat room
		chatRoom.AddClient(conn)

		// Read messages from the client and send to broadcast channel
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}
			chatRoom.broadcast <- msg
		}
	})

	fmt.Println("Chat server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
