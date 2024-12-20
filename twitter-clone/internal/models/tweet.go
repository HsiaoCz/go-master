package models

type Tweet struct {
	Base
	Content  string    `gorm:"not null" json:"content"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	User     User      `json:"user"`
	Likes    []Like    `json:"likes"`
	Retweets []Retweet `json:"retweets"`
	ParentID *uint     `json:"parent_id,omitempty"` // For replies
	Parent   *Tweet    `json:"parent,omitempty"`
}
