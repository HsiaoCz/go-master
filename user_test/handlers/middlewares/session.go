package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/user_test/db"
	"github.com/HsiaoCz/go-master/user_test/types"
)

func SessionMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is logged in
		// For now, we'll just check if the session cookie exists
		// In production, implement proper session validation
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "session cookie required", http.StatusUnauthorized)
			return
		}

		// Get the session from the database
		var session types.Sessions
		err = db.Get().NewSelect().Model(&session).Where("session_id = ?", sessionCookie.Value).Scan(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the session has expired
		// For now, we'll just check if the session exists
		// In production, implement proper session expiration
		if time.Now().After(session.ExpiresAt) {
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}
		// Add user information to the request context
		ctx := context.WithValue(r.Context(), types.CtxSessionKey, session)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
