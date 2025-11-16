package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/pmollerus23/go-aws-server/internal/auth"
	"github.com/pmollerus23/go-aws-server/internal/aws"
	"github.com/pmollerus23/go-aws-server/internal/config"
	"github.com/pmollerus23/go-aws-server/internal/middleware"
)

// Server represents the HTTP server.
type Server struct {
	logger      *slog.Logger
	config      *config.Config
	awsClients  *aws.Clients
	authService *auth.CognitoService
	httpServer  *http.Server
}

// New creates a new Server instance.
func New(logger *slog.Logger, cfg *config.Config, awsClients *aws.Clients) *Server {
	// Initialize Cognito authentication service
	authService := auth.NewCognitoService(awsClients.Cognito, cfg.Cognito, logger)

	return &Server{
		logger:      logger,
		config:      cfg,
		awsClients:  awsClients,
		authService: authService,
	}
}

// Run starts the HTTP server and handles graceful shutdown.
func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Create HTTP handler
	handler := s.setupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         net.JoinHostPort(s.config.Server.Host, s.config.Server.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second, // Time to read request headers and body
		WriteTimeout: 15 * time.Second, // Time to write response
		IdleTimeout:  60 * time.Second, // Time to keep connection alive when idle
	}

	// Start server in goroutine
	go func() {
		s.logger.Info("server starting", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for shutdown signal
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.logger.Info("server shutting down")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

// setupRoutes configures all routes and middleware.
func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	s.registerRoutes(mux)

	// Apply middleware in reverse order (last one wraps all others)
	var handler http.Handler = mux
	handler = middleware.Logging(s.logger)(handler)
	handler = middleware.RequestSizeLimit(10 * 1024 * 1024)(handler) // 10MB limit
	handler = middleware.PanicRecovery(s.logger)(handler)

	return handler
}
