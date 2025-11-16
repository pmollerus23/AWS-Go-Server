package middleware

import (
	"net/http"
)

// RequestSizeLimit creates a middleware that limits request body size.
func RequestSizeLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			h.ServeHTTP(w, r)
		})
	}
}
