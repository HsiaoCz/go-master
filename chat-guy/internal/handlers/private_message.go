package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
	"chat-guy/internal/models"

	"github.com/gorilla/mux"
)

type privateMessageResponse struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	Username   string    `json:"username"`
	AvatarURL  string    `json:"avatar_url"`
	Read       bool      `json:"read"`
	CreatedAt  time.Time `json:"created_at"`
}

func SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	vars := mux.Vars(r)
	receiverID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	var message struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	// Check if they are friends
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM friends 
			WHERE user_id = $1 AND friend_id = $2 AND status = 'accepted'
		)`, claims.UserID, receiverID).Scan(&exists)

	if err != nil || !exists {
		http.Error(w, "Not friends with this user", http.StatusForbidden)
		return
	}

	var pm models.PrivateMessage
	err = db.QueryRow(`
		INSERT INTO private_messages (sender_id, receiver_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		claims.UserID, receiverID, message.Content, time.Now()).Scan(&pm.ID, &pm.CreatedAt)

	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	// Get sender's avatar URL
	var senderAvatar string
	err = db.QueryRow("SELECT avatar_url FROM users WHERE id = $1", claims.UserID).Scan(&senderAvatar)
	if err != nil {
		senderAvatar = "" // Use default avatar if not found
	}

	// Create response with additional fields
	response := privateMessageResponse{
		ID:         pm.ID,
		SenderID:   claims.UserID,
		ReceiverID: receiverID,
		Content:    message.Content,
		Username:   claims.Username,
		AvatarURL:  senderAvatar,
		CreatedAt:  pm.CreatedAt,
	}

	// Notify the receiver through WebSocket if they're online
	privateClientsMutex.RLock()
	if client, ok := privateClients[receiverID]; ok {
		messageJSON, err := json.Marshal(response)
		if err == nil {
			select {
			case client.send <- messageJSON:
			default:
				// Channel is full or closed, skip sending
			}
		}
	}
	privateClientsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	vars := mux.Vars(r)
	otherUserID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	rows, err := db.Query(`
		SELECT 
			pm.id, 
			pm.sender_id, 
			pm.receiver_id, 
			pm.content, 
			pm.read, 
			pm.created_at,
			u.username,
			u.avatar_url
		FROM private_messages pm
		JOIN users u ON pm.sender_id = u.id
		WHERE (pm.sender_id = $1 AND pm.receiver_id = $2)
		   OR (pm.sender_id = $2 AND pm.receiver_id = $1)
		ORDER BY pm.created_at DESC
		LIMIT 50`,
		claims.UserID, otherUserID)

	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []privateMessageResponse
	for rows.Next() {
		var msg privateMessageResponse
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.Read,
			&msg.CreatedAt,
			&msg.Username,
			&msg.AvatarURL,
		)
		if err != nil {
			http.Error(w, "Error scanning messages", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	// Mark messages as read
	_, err = db.Exec(`
		UPDATE private_messages
		SET read = true
		WHERE receiver_id = $1 AND sender_id = $2 AND read = false`,
		claims.UserID, otherUserID)

	if err != nil {
		http.Error(w, "Error marking messages as read", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
