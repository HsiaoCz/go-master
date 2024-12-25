package types

import (
	"time"

	"github.com/google/uuid"
)

type Contacts struct {
	ContactID     int64     `json:"contact_id"`
	UserID        uuid.UUID `json:"user_id"`
	ContactUserID uuid.UUID `json:"contact_user_id"`
	IsBlocked     bool      `json:"is_blocked"`
	AddedAt       time.Time `json:"added_at"`
}
