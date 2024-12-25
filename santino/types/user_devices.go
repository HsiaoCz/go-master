package types

import (
	"time"

	"github.com/google/uuid"
)

type UserDevices struct {
	DeviceID    int64     `json:"device_id"`
	UserID      uuid.UUID `json:"user_id"`
	DeviceName  string    `json:"device_name"`
	DeviceToken string    `json:"device_token"`
	LastLogin   time.Time `json:"last_login"`
}
