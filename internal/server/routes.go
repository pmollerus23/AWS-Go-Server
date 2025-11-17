package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pmollerus23/go-aws-server/internal/handlers"
	"github.com/pmollerus23/go-aws-server/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// registerRoutes registers all HTTP routes.
func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Health check (public)
	mux.HandleFunc("GET /healthz", handlers.HandleHealthz(s.logger))

	// Auth endpoints (public)
	mux.Handle("POST /api/v1/auth/signup", handlers.HandleSignUp(s.logger, s.authService))
	mux.Handle("POST /api/v1/auth/confirm", handlers.HandleConfirmSignUp(s.logger, s.authService))
	mux.Handle("POST /api/v1/auth/login", handlers.HandleLogin(s.logger, s.authService))
	mux.Handle("POST /api/v1/auth/refresh", handlers.HandleRefreshToken(s.logger, s.authService))
	mux.Handle("POST /api/v1/auth/forgot-password", handlers.HandleForgotPassword(s.logger, s.authService))
	mux.Handle("POST /api/v1/auth/reset-password", handlers.HandleConfirmForgotPassword(s.logger, s.authService))

	// Protected routes - apply authentication middleware
	authMiddleware := middleware.Authenticate(s.authService, s.logger)

	// Item CRUD operations (protected)
	mux.Handle("GET /api/v1/items", authMiddleware(handlers.HandleItemsGet(s.logger)))
	mux.Handle("POST /api/v1/items", authMiddleware(handlers.HandleItemsCreate(s.logger)))

	// AWS service endpoints (protected)
	mux.Handle("GET /api/v1/aws/s3/buckets", authMiddleware(handlers.HandleS3ListBuckets(s.logger, s.awsClients.S3)))
	mux.Handle("GET /api/v1/aws/dynamodb/tables", authMiddleware(handlers.HandleDynamoDBListTables(s.logger, s.awsClients.DynamoDB)))
	mux.Handle("GET /api/v1/aws/dynamodb/records", authMiddleware(handlers.HandleDynamoDBListRecords(s.logger, s.awsClients.DynamoDB)))
	mux.Handle("POST /api/v1/aws/dynamodb/tables", authMiddleware(handlers.HandleDynamoDBUpsertTable(s.logger, s.awsClients.DynamoDB)))

	// Swagger documentation (public)
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", httpSwagger.WrapHandler))

	// Serve static files from React app (must be last to act as fallback)
	mux.Handle("/", s.spaHandler())
}

// spaHandler serves the React SPA from web/dist directory.
// It handles client-side routing by serving index.html for routes that don't exist.
func (s *Server) spaHandler() http.Handler {
	spaDir := "web/dist"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Build the path to the requested file
		path := filepath.Join(spaDir, r.URL.Path)

		// Check if file exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// File does not exist, serve index.html for client-side routing
			http.ServeFile(w, r, filepath.Join(spaDir, "index.html"))
			return
		} else if err != nil {
			// Error checking file
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Serve the requested file
		http.FileServer(http.Dir(spaDir)).ServeHTTP(w, r)
	})
}
