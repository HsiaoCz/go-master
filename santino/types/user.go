package types

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	UserID         uuid.UUID `gorm:"primaryKey" json:"user_id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	PasswordHash   string    `json:"-"`
	ProfilePicture string    `json:"profile_picture"`
	StatusMessage  string    `json:"status_message"`
	OnlineStatus   string    `json:"online_status"`
	LastSeen       time.Time `json:"last_seen"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
