package types

import (
	"time"

	"github.com/google/uuid"
)

type EncryptionKeys struct {
	KeyID      int64     `json:"key_id"`
	UserID     uuid.UUID `json:"user_id"`
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private"`
	CreatedAt  time.Time `json:"created_at"`
}
