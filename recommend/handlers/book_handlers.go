package handlers

import (
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
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
	return nil
}
