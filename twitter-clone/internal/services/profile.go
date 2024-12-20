package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"twitter-clone/internal/database"
	"twitter-clone/internal/models"

	"math/rand"
)

type ProfileService struct {
	uploadDir string
}

func NewProfileService() *ProfileService {
	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/profiles"
	os.MkdirAll(uploadDir, 0755)
	return &ProfileService{uploadDir: uploadDir}
}

func (s *ProfileService) UpdateBackground(userID uint, file *multipart.FileHeader) (string, error) {
	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return "", errors.New("invalid file type. Only jpeg, png, and gif are allowed")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("bg_%d_%s%s", userID, generateRandomString(8), ext)
	filepath := filepath.Join(s.uploadDir, filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Update user's background URL in database
	backgroundURL := fmt.Sprintf("/uploads/profiles/%s", filename)
	err = database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("background_url", backgroundURL).Error

	if err != nil {
		// Clean up file if database update fails
		os.Remove(filepath)
		return "", err
	}

	return backgroundURL, nil
}

func (s *ProfileService) UpdateProfile(userID uint, bio string) error {
	return database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("bio", bio).Error
}

func isValidImageType(contentType string) bool {
	validTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, t := range validTypes {
		if t == contentType {
			return true
		}
	}
	return false
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
