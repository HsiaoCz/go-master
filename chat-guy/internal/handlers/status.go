package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
	"chat-guy/internal/models"

	"github.com/gorilla/mux"
)

const (
	statusMediaDir = "./uploads/status"
	maxStatusSize  = 10 << 20 // 10MB
)

func init() {
	if err := os.MkdirAll(statusMediaDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create status media directory: %v", err))
	}
}

func CreateStatus(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	// Parse multipart form for media upload
	if err := r.ParseMultipartForm(maxStatusSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	var mediaURL string
	file, header, err := r.FormFile("media")
	if err == nil {
		defer file.Close()

		// Create unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("status_%d_%d%s", claims.UserID, time.Now().Unix(), ext)
		filepath := filepath.Join(statusMediaDir, filename)

		dst, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "Error saving media", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving media", http.StatusInternalServerError)
			return
		}

		mediaURL = fmt.Sprintf("/uploads/status/%s", filename)
	}

	// Set expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	db := database.GetDB()
	var status models.StatusUpdate
	err = db.QueryRow(`
		INSERT INTO status_updates (user_id, content, media_url, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`,
		claims.UserID, content, mediaURL, expiresAt).Scan(
		&status.ID,
		&status.CreatedAt,
		&status.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "Error creating status", http.StatusInternalServerError)
		return
	}

	status.UserID = claims.UserID
	status.Content = content
	status.MediaURL = mediaURL
	status.ExpiresAt = &expiresAt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func GetFriendsStatuses(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	db := database.GetDB()
	rows, err := db.Query(`
		SELECT 
			s.id,
			s.user_id,
			s.content,
			s.media_url,
			s.expires_at,
			s.created_at,
			s.updated_at,
			u.username,
			u.avatar_url,
			COUNT(DISTINCT sv.viewer_id) as view_count,
			EXISTS(
				SELECT 1 FROM status_views 
				WHERE status_id = s.id AND viewer_id = $1
			) as viewed
		FROM status_updates s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN status_views sv ON s.id = sv.status_id
		WHERE s.user_id IN (
			SELECT friend_id FROM friends 
			WHERE user_id = $1 AND status = 'accepted'
		)
		AND s.expires_at > NOW()
		GROUP BY s.id, u.username, u.avatar_url
		ORDER BY s.created_at DESC`,
		claims.UserID)

	if err != nil {
		http.Error(w, "Error fetching statuses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var statuses []models.StatusUpdate
	for rows.Next() {
		var status models.StatusUpdate
		err := rows.Scan(
			&status.ID,
			&status.UserID,
			&status.Content,
			&status.MediaURL,
			&status.ExpiresAt,
			&status.CreatedAt,
			&status.UpdatedAt,
			&status.Username,
			&status.AvatarURL,
			&status.ViewCount,
			&status.Viewed,
		)
		if err != nil {
			http.Error(w, "Error scanning statuses", http.StatusInternalServerError)
			return
		}
		statuses = append(statuses, status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func ViewStatus(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)
	vars := mux.Vars(r)
	statusID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid status ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	_, err = db.Exec(`
		INSERT INTO status_views (status_id, viewer_id, viewed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (status_id, viewer_id) DO NOTHING`,
		statusID, claims.UserID)

	if err != nil {
		http.Error(w, "Error recording view", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
