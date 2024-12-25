package types

import (
	"time"

	"github.com/google/uuid"
)

type ChatThreads struct {
	ThreadID      uuid.UUID `json:"thread_id"`
	User1ID       uuid.UUID `json:"user1_id"`
	User2ID       uuid.UUID `json:"user2_id"`
	LastMessageID uuid.UUID `json:"last_message_id"`
	LastUpdated   time.Time `json:"last_updated"`
}
