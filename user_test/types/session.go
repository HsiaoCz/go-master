package types

import "time"

type Sessions struct {
	ID        int64     `json:"id"`
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	IpAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
