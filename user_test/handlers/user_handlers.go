package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/bunt/db"
	"github.com/HsiaoCz/go-master/user_test/types"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	// Register user
	var user_register types.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user_register); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := user_register.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &types.Users{
		Username:      user_register.Username,
		Email:         user_register.Email,
		Password:      user_register.Password,
		Bio:           "",
		Avatar:        "",
		BackgroundUrl: "",
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
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User created successfully", "status": http.StatusCreated})
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var user_login types.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user_login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exists
	// ...
	var user types.Users
	err := db.Get().NewSelect().Model(&user).Where("email = ?", user_login.Email).Scan(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if the password is correct
	if !user.CheckPassword(user_login.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	// ...
	// Generate a session
	
}
