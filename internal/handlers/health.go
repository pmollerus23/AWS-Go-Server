package handlers

import (
	"log/slog"
	"net/http"
)

// HandleHealthz returns a simple health check handler.
//
//	@Summary		Health Check
//	@Description	Check if the server is healthy and responding
//	@Tags			health
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Router			/healthz [get]
func HandleHealthz(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("health check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
