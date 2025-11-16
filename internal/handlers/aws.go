package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// HandleS3ListBuckets returns a handler that lists all S3 buckets.
//
//	@Summary		List S3 buckets
//	@Description	Get a list of all S3 buckets in the AWS account
//	@Tags			aws
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"buckets and count"
//	@Failure		401	{string}	string					"Unauthorized"
//	@Failure		500	{string}	string					"Failed to list S3 buckets"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets [get]
func HandleS3ListBuckets(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("listing S3 buckets")

		result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			logger.Error("failed to list S3 buckets", "error", err)
			http.Error(w, "Failed to list S3 buckets", http.StatusInternalServerError)
			return
		}

		// Convert to response format
		buckets := make([]map[string]interface{}, 0, len(result.Buckets))
		for _, bucket := range result.Buckets {
			buckets = append(buckets, map[string]interface{}{
				"name":         *bucket.Name,
				"creationDate": bucket.CreationDate,
			})
		}

		response := map[string]interface{}{
			"buckets": buckets,
			"count":   len(buckets),
		}

		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleDynamoDBListTables returns a handler that lists all DynamoDB tables.
//
//	@Summary		List DynamoDB tables
//	@Description	Get a list of all DynamoDB tables in the AWS account
//	@Tags			aws
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"tables and count"
//	@Failure		401	{string}	string					"Unauthorized"
//	@Failure		500	{string}	string					"Failed to list DynamoDB tables"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/dynamodb/tables [get]
func HandleDynamoDBListTables(logger *slog.Logger, dynamoDBClient *dynamodb.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("listing DynamoDB tables")

		result, err := dynamoDBClient.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
		if err != nil {
			logger.Error("failed to list DynamoDB tables", "error", err)
			http.Error(w, "Failed to list DynamoDB tables", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"tables": result.TableNames,
			"count":  len(result.TableNames),
		}

		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
