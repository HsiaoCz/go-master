package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/recommend/storage"
	"github.com/HsiaoCz/go-master/recommend/types"
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
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "no auther information,please login", http.StatusNonAuthoritativeInfo)
			return
		}
		session, err := s.sen.GetSessionByToken(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "no auther information,please login", http.StatusNonAuthoritativeInfo)
			return
		}
		if time.Now().After(session.ExpiresAt) {
			http.Error(w, "current login information expired,please login again", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), types.CtxUserInfoKey, session)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
