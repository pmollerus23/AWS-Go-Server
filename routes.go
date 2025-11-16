package main

import (
	"log/slog"
	"net/http"

	"github.com/pmollerus23/go-aws-server/handlers"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	config *Config,
	awsClients *AWSClients,
) {
	mux.Handle("GET /api/v1/items", handlers.HandleItemsGet(logger))
	mux.Handle("POST /api/v1/items", handlers.HandleItemsCreate(logger))
	mux.HandleFunc("GET /healthz", handlers.HandleHealthz(logger))

	// AWS-enabled handlers
	mux.Handle("GET /api/v1/aws/s3/buckets", handlers.HandleS3ListBuckets(logger, awsClients.S3))
	mux.Handle("GET /api/v1/aws/dynamodb/tables", handlers.HandleDynamoDBListTables(logger, awsClients.DynamoDB))

	mux.Handle("/", http.NotFoundHandler())
}
