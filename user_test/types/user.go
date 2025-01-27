package types

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyUsername = errors.New("empty username")

	ErrEmptyPassword = errors.New("empty password")

	ErrEmptyEmail = errors.New("empty email")

	ErrInvalidEmail = errors.New("invalid email")

	ErrPasswordNotMatch = errors.New("password does not match")
)

type Users struct {
	ID            int64     `json:"id"`
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	Avatar        string    `json:"avatar"`
	Bio           string    `json:"bio"`
	BackgroundUrl string    `json:"background_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserRegister struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
	Email      string `json:"email"`
}

func (u *UserRegister) Validate() error {
	if u.Username == "" {
		return ErrEmptyUsername
	}
	if u.Password == "" {
		return ErrEmptyPassword
	}
	if u.RePassword != u.Password {
		return ErrPasswordNotMatch
	}
	if u.Email == "" {
		return ErrEmptyEmail
	}
	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

func isValidEmail(email string) bool {
	// 正则表达式用于验证电子邮件地址
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func (u *Users) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *Users) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
