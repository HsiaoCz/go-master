package types

import (
	"time"

	"github.com/google/uuid"
)

type ShareReactions struct {
	ReactionID int64     `json:"reaction_id"`
	ShareID    uuid.UUID `json:"share_id"`
	UserID     uuid.UUID `json:"user_id"`
	// Reaction 'like' or 'dislike' or 'love' or 'haha' or 'wow' or 'sad' or 'angry'
	Reaction  string    `json:"reaction"`
	ReactedAt time.Time `json:"reacted_at"`
}
