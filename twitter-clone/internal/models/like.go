package models

type Like struct {
	Base
	UserID  uint  `gorm:"not null" json:"user_id"`
	TweetID uint  `gorm:"not null" json:"tweet_id"`
	User    User  `json:"user"`
	Tweet   Tweet `json:"tweet"`
}