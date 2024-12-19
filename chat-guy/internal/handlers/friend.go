package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
	"chat-guy/internal/models"

	"github.com/gorilla/mux"
)

func AddFriend(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	vars := mux.Vars(r)
	friendUsername := vars["username"]

	db := database.GetDB()
	var friendID int64
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", friendUsername).Scan(&friendID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if friendID == claims.UserID {
		http.Error(w, "Cannot add yourself as friend", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		INSERT INTO friends (user_id, friend_id, status, created_at, updated_at)
		VALUES ($1, $2, 'pending', $3, $3)
		ON CONFLICT (user_id, friend_id) DO NOTHING`,
		claims.UserID, friendID, time.Now())

	if err != nil {
		http.Error(w, "Error adding friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	vars := mux.Vars(r)
	friendID := vars["id"]

	db := database.GetDB()
	_, err := db.Exec(`
		UPDATE friends 
		SET status = 'accepted', updated_at = $3
		WHERE user_id = $1 AND friend_id = $2 AND status = 'pending'`,
		friendID, claims.UserID, time.Now())

	if err != nil {
		http.Error(w, "Error accepting friend request", http.StatusInternalServerError)
		return
	}

	// Create reverse friendship
	_, err = db.Exec(`
		INSERT INTO friends (user_id, friend_id, status, created_at, updated_at)
		VALUES ($1, $2, 'accepted', $3, $3)`,
		claims.UserID, friendID, time.Now())

	if err != nil {
		http.Error(w, "Error creating reverse friendship", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	db := database.GetDB()
	rows, err := db.Query(`
		SELECT f.friend_id, u.username, f.status, f.created_at, f.updated_at
		FROM friends f
		JOIN users u ON f.friend_id = u.id
		WHERE f.user_id = $1`,
		claims.UserID)

	if err != nil {
		http.Error(w, "Error fetching friends", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var friends []models.Friend
	for rows.Next() {
		var friend models.Friend
		err := rows.Scan(
			&friend.FriendID,
			&friend.Username,
			&friend.Status,
			&friend.CreatedAt,
			&friend.UpdatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning friends", http.StatusInternalServerError)
			return
		}
		friends = append(friends, friend)
	}

	json.NewEncoder(w).Encode(friends)
}

func GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	db := database.GetDB()
	rows, err := db.Query(`
		SELECT f.user_id, u.username, f.created_at
		FROM friends f
		JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = $1 AND f.status = 'pending'`,
		claims.UserID)

	if err != nil {
		http.Error(w, "Error fetching friend requests", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []models.Friend
	for rows.Next() {
		var request models.Friend
		err := rows.Scan(
			&request.UserID,
			&request.Username,
			&request.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning friend requests", http.StatusInternalServerError)
			return
		}
		requests = append(requests, request)
	}

	json.NewEncoder(w).Encode(requests)
}
