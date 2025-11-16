package server

import (
	"net/http"

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

	// 404 handler
	mux.Handle("/", http.NotFoundHandler())
}
