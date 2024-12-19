package handlers

import (
	"encoding/json"
	"net/http"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
	"chat-guy/internal/models"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	db := database.GetDB()
	var user models.User
	err := db.QueryRow(`
		SELECT id, username, email, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1`,
		claims.UserID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
