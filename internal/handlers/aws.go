package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pmollerus23/go-aws-server/internal/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	// "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

// HandleDynamoDBListRecords returns a handler that lists all records from a DynamoDB table.
//
//	@Summary		List DynamoDB records
//	@Description	Get a list of all records from a DynamoDB table
//	@Tags			aws
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"records and count"
//	@Failure		401	{string}	string					"Unauthorized"
//	@Failure		500	{string}	string					"Failed to list records"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/dynamodb/records [get]
func HandleDynamoDBListRecords(logger *slog.Logger, dynamoDBClient *dynamodb.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Listing records from DynamoDB table")

		tableName := "Phil_Go_App_Database"
		result, err := dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
			TableName: aws.String(tableName),
		})

		if err != nil {
			logger.Error("Failed to scan DynamoDB table", "error", err, "table", tableName)
			http.Error(w, "Failed to list records", http.StatusInternalServerError)
			return
		}

		// Unmarshal the items into our model
		var records []models.DynamoDBRecord
		err = attributevalue.UnmarshalListOfMaps(result.Items, &records)
		if err != nil {
			logger.Error("Failed to unmarshal DynamoDB items", "error", err)
			http.Error(w, "Failed to process records", http.StatusInternalServerError)
			return
		}

		logger.Info("Successfully retrieved records", "count", len(records))

		response := map[string]interface{}{
			"records": records,
			"count":   len(records),
		}

		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleDynamoDBUpsertTable returns a handler that inserts or updates a record in a DynamoDB table.
//
//	@Summary		Upsert DynamoDB record
//	@Description	Insert or update a record in a DynamoDB table
//	@Tags			aws
//	@Accept			json
//	@Produce		json
//	@Param			record	body		models.DynamoDBRecord		true	"Record to upsert"
//	@Success		201		{object}	map[string]interface{}		"result metadata"
//	@Failure		400		{string}	string						"Invalid request body"
//	@Failure		401		{string}	string						"Unauthorized"
//	@Failure		500		{string}	string						"Failed to upsert record"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/dynamodb/tables [post]
func HandleDynamoDBUpsertTable(logger *slog.Logger, dynamoDBClient *dynamodb.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Upserting record into DynamoDB table")

		// Decode the JSON payload from the request body
		var record models.DynamoDBRecord
		if err := decode(r, &record); err != nil {
			logger.Error("Failed to decode request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		logger.Info("Decoded record", "id", record.ID, "name", record.Name, "updated_at", record.UpdatedAt)

		item, err := attributevalue.MarshalMap(record)
		if err != nil {
			logger.Error("Failed to marshal user request record into DynamoDB object", "error", err)
			http.Error(w, "Failed to marshal user request record into DynamoDB object", http.StatusInternalServerError)
			return
		}

		logger.Info("Marshaled item", "item", item)

		tableName := "Phil_Go_App_Database"
		logger.Info("Putting item to DynamoDB", "table", tableName)

		result, err := dynamoDBClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})

		if err != nil {
			logger.Error("Failed to put record in DynamoDB", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Info("Successfully put item to DynamoDB", "result", result)

		response := map[string]interface{}{
			"result_attributes": result.Attributes,
			"success":           true,
		}

		if err := encode(w, r, int(http.StatusCreated), response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}
