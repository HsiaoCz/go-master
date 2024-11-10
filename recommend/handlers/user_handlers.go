package handlers

import (
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
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
	return nil
}
