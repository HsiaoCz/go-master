package middleware

import (
	"net/http"
	"strconv"
	"time"

	"chat-guy/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		status := wrapped.status
		if status == 0 {
			status = 200
		}

		metrics.HTTPRequestDuration.WithLabelValues(
			r.URL.Path,
			r.Method,
			strconv.Itoa(status),
		).Observe(time.Since(start).Seconds())
	})
}
