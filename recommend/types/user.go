package types

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	UserID       string `gorm:"column:user_id;" json:"user_id"`
	Phone        string `gorm:"column:phone;" json:"phone"`
	HashPassword string `gorm:"column:hash_password;" json:"hash_password"`
	Username     string `gorm:"column:username;" json:"username"`
	Role         bool   `gorm:"column:role;" json:"role"`
}
