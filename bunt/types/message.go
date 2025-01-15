package types

import (
	"time"
)

type Messages struct {
	ID         int64     `json:"id"`
	MessageID  string `json:"message_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	Type       string    `json:"type"` // "text","image","system"
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"create_at"`
}
