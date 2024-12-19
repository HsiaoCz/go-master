package models

import "time"

type Friend struct {
	UserID    int64     `json:"user_id"`
	FriendID  int64     `json:"friend_id"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PrivateMessage struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Content    string    `json:"content"`
	Read       bool      `json:"read"`
	CreatedAt  time.Time `json:"created_at"`
}
