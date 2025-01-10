package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/bunt/db"
	"github.com/HsiaoCz/go-master/bunt/types"
)

func SessionMiddleware(next http.Handler) http.HandlerFunc {
	// Check if the user is logged in
	// ...
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		var session types.Sessions
		err = db.Get().NewSelect().Model(&session).Where("session_id = ?", cookie.Value).Scan(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if time.Now().After(session.ExpiresAt) {
			http.Error(w, "please login", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), types.CtxSessionKey, &session)
		r = r.WithContext(ctx)
		
		next.ServeHTTP(w, r)
	}

}
