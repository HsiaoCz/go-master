package types

import "time"

type MessageStatus struct {
	ID        int64     `json:"id"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	IsRead    bool      `json:"is_read"`
	ReadAt    time.Time `json:"read_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
