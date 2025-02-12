package types

import (
	"time"

	"github.com/google/uuid"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrUsernameEmpty = &Error{
		Code:    400,
		Message: "Username cannot be empty",
	}
	ErrEmailEmpty = &Error{
		Code:    400,
		Message: "Email cannot be empty",
	}

	ErrPasswordEmpty = &Error{
		Code:    400,
		Message: "Password cannot be empty",
	}
)

type Users struct {
	UserID         uuid.UUID `gorm:"primaryKey" json:"user_id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phone_number"`
	PasswordHash   string    `json:"-"`
	ProfilePicture string    `json:"profile_picture"`
	StatusMessage  string    `json:"status_message"`
	OnlineStatus   string    `json:"online_status"`
	LastSeen       time.Time `json:"last_seen"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRegister struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	Phone      string `json:"phone_number"`
	Avatar     string `json:"avatar"`
	Background string `json:"background_url"`
}

func (u *UserRegister) Validate() error {
	if u.Username == "" {
		return ErrUsernameEmpty
	}
	if u.Email == "" {
		return ErrEmailEmpty
	}
	if u.Password == "" {
		return ErrPasswordEmpty
	}
	return nil
}

func (u *Users) HashPassword() error {
	hash, err := HashPassword(u.PasswordHash)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

func HashPassword(password string) (string, error) {
	return password, nil
}
