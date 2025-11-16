package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// PanicRecovery creates a middleware that recovers from panics.
func PanicRecovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
						"remote_addr", r.RemoteAddr,
						"stack", string(debug.Stack()),
					)

					// Send 500 error to client
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"Internal server error"}`)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
