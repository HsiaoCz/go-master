package types

import (
	"time"

	"github.com/google/uuid"
)

type ShareComments struct {
	CommentID   uuid.UUID `json:"comment_id"`
	ShareID     uuid.UUID `json:"share_id"`
	UserID      uuid.UUID `json:"user_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}
