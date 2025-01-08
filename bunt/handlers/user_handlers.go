package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/bunt/db"
	"github.com/HsiaoCz/go-master/bunt/types"
	"github.com/google/uuid"
)

func HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	var user_signup types.UserSignup
	if err := json.NewDecoder(r.Body).Decode(&user_signup); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user_signup.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Save the user to the database
	// ...
	user := &types.Users{
		Username:      user_signup.Username,
		Email:         user_signup.Email,
		PasswordHash:  user_signup.Password,
		UserID:        uuid.New().String(),
		BackgroundUrl: "",
		Bio:           "",
		AvatarUrl:     "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := user.HashPassword(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the user to the database
	// ...
	_, err := db.Get().NewInsert().Model(user).Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"message": "User created successfully", "status": http.StatusCreated})
}

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var user_login types.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user_login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user_login.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Find the user in the database
	// ...
	var user types.Users
	err := db.Get().NewSelect().Model(&user).Where("email = ?", user_login.Email).Scan(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.CheckPassword(user_login.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "User logged in successfully", "status": http.StatusOK})
}
