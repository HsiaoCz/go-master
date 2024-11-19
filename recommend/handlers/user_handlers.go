package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
	"github.com/HsiaoCz/go-master/recommend/types"
)

type UserHandlers struct {
	mod *mod.UserMod
}

func UserHandlersInit(mod *mod.UserMod) *UserHandlers {
	return &UserHandlers{
		mod: mod,
	}
}

func (u *UserHandlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var create_user_params types.Users
	if err := json.NewDecoder(r.Body).Decode(&create_user_params); err != nil {
		return ErrorMessage(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (u *UserHandlers) HandleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	user_id := r.PathValue("user_id")
	return WriteJson(w, http.StatusOK, user_id)
}
