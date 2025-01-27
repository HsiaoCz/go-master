package handlers

import (
	"encoding/json"
	"net/http"

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
	// Save the user to the database
	// ...
}
