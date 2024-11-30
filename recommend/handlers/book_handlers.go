package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
	"github.com/HsiaoCz/go-master/recommend/types"
)

type BookHandlers struct {
	book mod.BookModInter
}

func BookHandlersInit(book mod.BookModInter) *BookHandlers {
	return &BookHandlers{
		book: book,
	}
}

func (b *BookHandlers) HandleCreateBook(w http.ResponseWriter, r *http.Request) error {
	userInfo, ok := r.Context().Value(types.CtxUserInfoKey).(*types.UserInfo)
	if !ok {
		return ErrorMessage(http.StatusNonAuthoritativeInfo, "please login")
	}
	if !userInfo.Role {
		return ErrorMessage(http.StatusNotAcceptable, "need admin")
	}
	var book types.CreateBookParams
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		return ErrorMessage(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (b *BookHandlers) HandleGetBookByAuther(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *BookHandlers) HandleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *BookHandlers) HandleGetBookByRecords(w http.ResponseWriter, r *http.Request) error {
	return nil
}
