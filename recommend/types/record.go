package types

import "gorm.io/gorm"

type Records struct {
	gorm.Model
	UserID string `gorm:"column:user_id;" json:"user_id"`
	BookID string 
	Type   string
	Device string
}
