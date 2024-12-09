package middlewares

import (
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/storage"
)

type SessionAuth struct {
	sen storage.SessionStorer
}

func SessionAuthInit(sen storage.SessionStorer) *SessionAuth {
	return &SessionAuth{
		sen: sen,
	}
}

func (s *SessionAuth) AuthMiddlewares(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
