package types

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	MessageID  uuid.UUID `json:"message_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	IsRead     bool      `json:"is_read"`
	IsDeleted  bool      `json:"is_deleted"`
}
