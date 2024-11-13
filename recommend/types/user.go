package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	UserID          string `gorm:"column:user_id;" json:"user_id"`
	Phone           string `gorm:"column:phone;" json:"phone"`
	HashPassword    string `gorm:"column:hash_password;" json:"-"`
	Username        string `gorm:"column:username;" json:"username"`
	Role            bool   `gorm:"column:role;" json:"role"`
	Avatar          string `gorm:"column:avatar;" json:"avatar"`
	Brief           string `gorm:"column:brief;" json:"brief"`
	Birthday        string `gorm:"column:birthday;" json:"birthday"`
	Age             int    `gorm:"column:age;" json:"age"`
	BackgroundImage string `gorm:"column:background_image;" json:"background_image"`
	Gender          string `gorm:"column:gender;" json:"gender"`
}

type CreateUserParams struct {
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     bool   `json:"role"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
}

func CreateUserFromParams(params CreateUserParams) *Users {
	return &Users{
		UserID:          uuid.NewString(),
		Phone:           params.Phone,
		Username:        params.Username,
		Role:            params.Role,
		Birthday:        params.Birthday,
		Gender:          params.Gender,
		Avatar:          "./picture/avatar/1233.jpg",
		Brief:           "",
		BackgroundImage: "./picture/background/1344.jpg",
		Age:             27,
	}
}
