package types

import "time"

type Posts struct {
	ID        int64     `json:"id" `
	PostID    string    `json:"post_id" `
	UserID    string    `json:"user_id"`
	Caption   string    `json:"caption"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url"`
	VideoUrl  string    `json:"video_url"`
	Location  string    `json:"location"`
	Likes     int64     `json:"likes"`
	Comments  int64     `json:"comments"`
	CreatedAt time.Time `json:"created_at" `
	UpdatedAt time.Time `json:"updated_at"`
}
