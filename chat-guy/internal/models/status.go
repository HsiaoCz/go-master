package models

import "time"

type StatusUpdate struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Content   string     `json:"content"`
	MediaURL  string     `json:"media_url,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	
	// Additional fields for response
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	ViewCount int    `json:"view_count,omitempty"`
	Viewed    bool   `json:"viewed,omitempty"`
}

type StatusView struct {
	StatusID int64     `json:"status_id"`
	ViewerID int64     `json:"viewer_id"`
	ViewedAt time.Time `json:"viewed_at"`
}
