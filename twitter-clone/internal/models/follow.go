package models

type Follow struct {
	Base
	FollowerID  uint `gorm:"not null" json:"follower_id"`
	FollowingID uint `gorm:"not null" json:"following_id"`
	Follower    User `json:"follower"`
	Following   User `json:"following"`
}