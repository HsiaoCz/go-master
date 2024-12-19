package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"chat-guy/internal/database"
	"chat-guy/internal/middleware"
)

const (
	maxUploadSize = 5 << 20 // 5MB
	uploadDir     = "./uploads/avatars"
)

func init() {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*middleware.Claims)

	// Parse multipart form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isValidImageType(header.Filename) {
		http.Error(w, "Invalid file type. Only jpg, jpeg, png allowed", http.StatusBadRequest)
		return
	}

	// Create unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", claims.UserID, ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create new file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Update user's avatar URL in database
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	db := database.GetDB()
	_, err = db.Exec(`
		UPDATE users 
		SET avatar_url = $1, updated_at = NOW()
		WHERE id = $2`,
		avatarURL, claims.UserID)

	if err != nil {
		http.Error(w, "Error updating avatar URL", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"avatar_url": "%s"}`, avatarURL)
}

func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func ServeAvatar(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	filepath := filepath.Join(uploadDir, filename)

	// Serve the file
	http.ServeFile(w, r, filepath)
}
