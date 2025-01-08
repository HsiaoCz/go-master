package types

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
)

type Users struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	UserID        string    `json:"user_id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"password_hash"`
	AvatarUrl     string    `json:"avatar_url"`
	BackgroundUrl string    `json:"background_url"`
	Bio           string    `json:"bio"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserSignup struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserSignup) Validate() error {
	regexpEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexpEmail.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}
func (u *Users) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func (u *Users) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserLogin) Validate() error {
	regexpEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexpEmail.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}
