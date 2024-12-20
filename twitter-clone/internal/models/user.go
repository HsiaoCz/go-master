package models

type User struct {
	Base
	Username        string    `gorm:"uniqueIndex;not null" json:"username"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	Password        string    `gorm:"not null" json:"-"`
	Bio             string    `json:"bio"`
	Avatar          string    `json:"avatar"`
	BackgroundImage string    `json:"background_image"`
	Followers       []Follow  `gorm:"foreignKey:FollowingID" json:"followers"`
	Following       []Follow  `gorm:"foreignKey:FollowerID" json:"following"`
	Tweets          []Tweet   `json:"tweets"`
	Likes           []Like    `json:"likes"`
	Retweets        []Retweet `json:"retweets"`
}
