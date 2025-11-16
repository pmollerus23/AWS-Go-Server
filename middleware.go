package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// newLoggingMiddleware creates a middleware that logs HTTP requests and responses.
func newLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info("request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)

			h.ServeHTTP(w, r)

			logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

// adminOnly is an example of a simple middleware without dependencies.
func adminOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: check if user is admin
		// In a real app, you would check authentication/authorization here
		isAdmin := r.Header.Get("X-Admin") == "true"
		if !isAdmin {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// newPanicRecoveryMiddleware creates a middleware that recovers from panics.
func newPanicRecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
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

// newRequestSizeLimitMiddleware creates a middleware that limits request body size.
func newRequestSizeLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			h.ServeHTTP(w, r)
		})
	}
}
