package models

import "time"

type Message struct {
	ID        int64     `json:"id"`
	RoomID    int64     `json:"room_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // "text", "image", "system"
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
	Private     bool      `json:"private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoomMember struct {
	RoomID    int64     `json:"room_id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"` // "admin", "member"
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
