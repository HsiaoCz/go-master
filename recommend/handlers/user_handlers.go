package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/recommend/storage"
	"github.com/HsiaoCz/go-master/recommend/types"
	"github.com/google/uuid"
)

type UserHandlers struct {
	mod storage.UserStorer
	sen storage.SessionStorer
	rec storage.RecordStorer
}

func UserHandlersInit(mod storage.UserStorer, sen storage.SessionStorer, rec storage.RecordStorer) *UserHandlers {
	return &UserHandlers{
		mod: mod,
		sen: sen,
		rec: rec,
	}
}

func (u *UserHandlers) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var create_user_params types.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&create_user_params); err != nil {
		return ErrorMessage(http.StatusBadRequest, err.Error())
	}
	user, err := u.mod.CreateUser(r.Context(), types.CreateUserFromParams(create_user_params))
	if err != nil {
		return ErrorMessage(http.StatusInternalServerError, err.Error())
	}
	return WriteJson(w, http.StatusOK, map[string]any{
		"status":  http.StatusOK,
		"message": "create user success",
		"user":    user,
	})
}

func (u *UserHandlers) HandleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	user_id := r.PathValue("user_id")
	user, err := u.mod.GetUserByID(r.Context(), user_id)
	if err != nil {
		return ErrorMessage(http.StatusBadRequest, err.Error())
	}
	return WriteJson(w, http.StatusOK, user)
}

func (u *UserHandlers) HandleDeleteUserByID(w http.ResponseWriter, r *http.Request) error {
	userInfo, ok := r.Context().Value(types.CtxUserInfoKey).(*types.UserInfo)
	if !ok {
		return ErrorMessage(http.StatusNonAuthoritativeInfo, "your have no rights to do this shit....")
	}
	if err := u.mod.DeleteUserByID(r.Context(), userInfo.UserID); err != nil {
		return ErrorMessage(http.StatusInternalServerError, err.Error())
	}
	return WriteJson(w, http.StatusOK, map[string]any{
		"status":  http.StatusOK,
		"message": "delete user success",
	})
}

func (u *UserHandlers) HandleUserLogin(w http.ResponseWriter, r *http.Request) error {
	var user_login_params types.UserLoginParams
	if err := json.NewDecoder(r.Body).Decode(&user_login_params); err != nil {
		return ErrorMessage(http.StatusBadRequest, err.Error())
	}
	user, err := u.mod.GetUserByPhoneAndPassword(r.Context(), &user_login_params)
	if err != nil {
		return ErrorMessage(http.StatusBadRequest, err.Error())
	}
	session := &types.Sessions{
		Token:     uuid.New().String(),
		UserID:    user.UserID,
		IpAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}
	session, err = u.sen.CreateSession(r.Context(), session)
	if err != nil {
		return ErrorMessage(http.StatusInternalServerError, err.Error())
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    session.Token,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
	})
	return WriteJson(w, http.StatusOK, map[string]any{
		"status":  http.StatusOK,
		"message": "login success",
		"user":    user,
	})
}

func (u *UserHandlers) HandleGetRecord(w http.ResponseWriter, r *http.Request) error {
	session, ok := r.Context().Value(types.CtxUserInfoKey).(*types.Sessions)
	if !ok {
		return ErrorMessage(http.StatusNonAuthoritativeInfo, "please login")
	}
	records, err := u.rec.GetRecordsByUserID(r.Context(), session.UserID)
	if err != nil {
		return ErrorMessage(http.StatusInternalServerError, "you have no record")
	}
	return WriteJson(w, http.StatusOK, records)
}
