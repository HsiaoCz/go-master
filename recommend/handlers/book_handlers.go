package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
	"github.com/HsiaoCz/go-master/recommend/types"
)

type BookHandlers struct {
	record mod.RecordModInter
	book   mod.BookModInter
}

func BookHandlersInit(book mod.BookModInter, record mod.RecordModInter) *BookHandlers {
	return &BookHandlers{
		book:   book,
		record: record,
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
	auther := r.URL.Query().Get("auther")
	books, err := b.book.GetBookByAuther(r.Context(), auther)
	if err != nil {
		return ErrorMessage(http.StatusInternalServerError, err.Error())
	}
	return WriteJson(w, http.StatusOK, map[string]any{
		"message": "get books success",
		"status":  http.StatusOK,
		"books":   books,
	})
}

func (b *BookHandlers) HandleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	book_id := r.PathValue("book_id")
	book, err := b.book.GetBookByID(r.Context(), book_id)
	if err != nil {
		return ErrorMessage(http.StatusInternalServerError, err.Error())
	}
	// after getting book we should create a record

	userInfo, ok := r.Context().Value(types.CtxUserInfoKey).(*types.UserInfo)
	if !ok {
		return ErrorMessage(http.StatusInternalServerError, "something wrong")
	}
	record_param := types.CreateRecordParams{
		BookID:     book_id,
		CoverImage: book.CoverImage,
		Title:      book.Title,
		Auther:     book.Auther,
		Device:     string(r.UserAgent()),
		UserID:     userInfo.UserID,
		TypeName:   "see",
	}
	record, err := b.record.CreateRecord(r.Context(), types.CreateRecordFromParams(record_param))
	if err != nil {
		return ErrorMessage(http.StatusInternalServerError, err.Error())
	}
	return WriteJson(w, http.StatusOK, map[string]any{
		"message": "get book success",
		"status":  http.StatusOK,
		"book":    book,
		"record":  record,
	})
}

func (b *BookHandlers) HandleGetBookByRecords(w http.ResponseWriter, r *http.Request) error {
	return nil
}
