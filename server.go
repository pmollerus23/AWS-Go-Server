package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func run(
	ctx context.Context,
	w io.Writer,
	args []string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Create logger
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Configuration (could be loaded from env or config file)
	config := &Config{
		Host: "localhost",
		Port: "8080",
		AWS: AWSConfig{
			Region:  getEnvOrDefault("AWS_REGION", "us-east-1"),
			Profile: getEnvOrDefault("AWS_PROFILE", ""),
		},
	}

	// Initialize AWS clients
	awsClients, err := NewAWSClients(ctx, logger, &config.AWS)
	if err != nil {
		return fmt.Errorf("failed to initialize AWS clients: %w", err)
	}

	// Create server
	srv := NewServer(logger, config, awsClients)

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      srv,
		ReadTimeout:  15 * time.Second, // Time to read request headers and body
		WriteTimeout: 15 * time.Second, // Time to write response
		IdleTimeout:  60 * time.Second, // Time to keep connection alive when idle
	}

	go func() {
		logger.Info("server starting", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Info("server shutting down")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

type Config struct {
	Host string
	Port string
	AWS  AWSConfig
}

// getEnvOrDefault returns the value of an environment variable or a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewServer(
	logger *slog.Logger,
	config *Config,
	awsClients *AWSClients,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		mux,
		logger,
		config,
		awsClients,
	)

	var handler http.Handler = mux
	// Apply middleware in reverse order (last one wraps all others)
	handler = newLoggingMiddleware(logger)(handler)
	handler = newRequestSizeLimitMiddleware(10 * 1024 * 1024)(handler) // 10MB limit
	handler = newPanicRecoveryMiddleware(logger)(handler)

	return handler
}
