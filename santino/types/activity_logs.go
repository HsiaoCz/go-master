package types

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type ActivityLogs struct {
	LogID             int64     `json:"log_id"`
	UserID            uuid.UUID `json:"user_id"`
	ActivityType      string    `json:"activity_type"`
	ActivityTimestamp time.Time `json:"activity_timestamp"`
	IpAddress         net.IP    `json:"ip_address"`
}
