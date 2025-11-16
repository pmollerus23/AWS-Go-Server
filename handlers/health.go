package handlers

import (
	"log/slog"
	"net/http"
)

// HandleHealthz returns a simple health check handler.
func HandleHealthz(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("health check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
