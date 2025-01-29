package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HsiaoCz/go-master/user_test/db"
	"github.com/HsiaoCz/go-master/user_test/types"
	"github.com/google/uuid"
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

	key := fmt.Sprintf("%s:login_count", user_login.Email)
	result := db.RDB.Get(r.Context(), key)
	count, err := strconv.Atoi(result.Val())
	if err != nil {
		count = 0
	}

	if count >= 5 {
		http.Error(w, "too many login attempts", http.StatusTooManyRequests)
		return
	}

	// Check if the user exists
	// ...
	var user types.Users
	err = db.Get().NewSelect().Model(&user).Where("email = ?", user_login.Email).Scan(r.Context())
	if err != nil {
		status := db.RDB.Set(r.Context(), key, count+1, time.Minute*5)
		if status.Err() != nil {
			http.Error(w, status.Err().Error(), http.StatusInternalServerError)
			return
		}
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

	session := &types.Sessions{
		UserID:    user.UserID,
		SessionID: uuid.New().String(),
		IpAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = db.Get().NewInsert().Model(session).Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the session
	// ...
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User logged in successfully", "status": http.StatusOK})
}
