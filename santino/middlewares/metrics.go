package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/HsiaoCz/go-master/santino/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}
	rw.statusCode = statusCode
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

// MetricsMiddleware is a middleware that records the duration of each HTTP request.
func MetricsMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		status := wrapped.statusCode
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

