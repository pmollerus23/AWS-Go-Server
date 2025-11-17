package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// HandleHealthz returns a simple health check handler.
//
//	@Summary		Health Check
//	@Description	Check if the server is healthy and responding
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Router			/healthz [get]
func HandleHealthz(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("health check")

		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("failed to encode health response", "error", err)
		}
	}
}
