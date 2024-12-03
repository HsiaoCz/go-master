package middlewares

import (
	"net/http"

	"github.com/HsiaoCz/go-master/recommend/mod"
)

type SessionAuth struct {
	sen mod.SessionModInter
}

func SessionAuthInit(sen mod.SessionModInter) *SessionAuth {
	return &SessionAuth{
		sen: sen,
	}
}

func (s *SessionAuth) AuthMiddlewares(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
