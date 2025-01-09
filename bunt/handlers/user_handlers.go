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

func HandleUserGet(w http.ResponseWriter, r *http.Request) {
	// Get the user from the database
	// ...
	var user types.Users
	err := db.Get().NewSelect().Model(&user).Where("user_id = ?", r.URL.Query().Get("user_id")).Scan(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"user": user, "status": http.StatusOK})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	var user_update types.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&user_update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Update the user in the database
	// ...
	_, err := db.Get().NewUpdate().Model(&types.Users{
		Username:      user_update.Username,
		BackgroundUrl: user_update.BackgroundUrl,
		Bio:           user_update.Bio,
		AvatarUrl:     user_update.AvatarUrl,
		UpdatedAt:     time.Now(),
	}).Where("user_id = ?", r.URL.Query().Get("user_id")).Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "User updated successfully", "status": http.StatusOK})
}

func HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	// Delete the user from the database
	// ...
	_, err := db.Get().NewDelete().Table("users").Where("user_id = ?", r.URL.Query().Get("user_id")).Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "User deleted successfully", "status": http.StatusOK})
}
