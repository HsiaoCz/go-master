package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/stock/storage"
	"github.com/HsiaoCz/go-master/stock/types"
)

type AuthSessionMiddleware struct {
	sen storage.SessionStoreInter
}

func AuthSessionMiddlewareInit(sen storage.SessionStoreInter) *AuthSessionMiddleware {
	return &AuthSessionMiddleware{
		sen: sen,
	}
}

func (a *AuthSessionMiddleware) AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "no auther information,please login", http.StatusNonAuthoritativeInfo)
			return
		}
		session, err := a.sen.GetSessionByToken(r.Context(), cookie.Value)
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
