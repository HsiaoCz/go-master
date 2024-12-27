package types

import (
	"time"

	"github.com/google/uuid"
)

type Shares struct {
	ShareID  uuid.UUID `json:"share_id"`
	UserID   uuid.UUID `json:"user_id"`
	Content  string    `json:"content"`
	MediaURL string    `json:"media_url"`
	// Visbility 'public' or 'private' or 'friends'
	Visbility string    `json:"visibility"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}
