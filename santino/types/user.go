package types

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
	
}
