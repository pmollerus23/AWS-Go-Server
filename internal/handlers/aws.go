package handlers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

// HandleS3CreateBucket creates a new S3 bucket.
//
//	@Summary		Create S3 bucket
//	@Description	Create a new S3 bucket
//	@Tags			aws
//	@Accept			json
//	@Produce		json
//	@Param			request	body		map[string]string	true	"Bucket name"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"Failed to create bucket"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets [post]
func HandleS3CreateBucket(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			BucketName string `json:"bucketName"`
			Region     string `json:"region"`
		}

		if err := decode(r, &req); err != nil {
			logger.Error("failed to decode request", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.BucketName == "" {
			http.Error(w, "Bucket name is required", http.StatusBadRequest)
			return
		}

		logger.Info("creating S3 bucket", "bucket", req.BucketName, "region", req.Region)

		input := &s3.CreateBucketInput{
			Bucket: aws.String(req.BucketName),
		}

		// If region is specified and not us-east-1, add location constraint
		if req.Region != "" && req.Region != "us-east-1" {
			input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(req.Region),
			}
		}

		_, err := s3Client.CreateBucket(context.TODO(), input)
		if err != nil {
			logger.Error("failed to create S3 bucket", "error", err)
			http.Error(w, fmt.Sprintf("Failed to create bucket: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":    true,
			"bucketName": req.BucketName,
		}

		if err := encode(w, r, http.StatusCreated, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleS3DeleteBucket deletes an S3 bucket.
//
//	@Summary		Delete S3 bucket
//	@Description	Delete an S3 bucket
//	@Tags			aws
//	@Produce		json
//	@Param			bucketName	path		string	true	"Bucket name"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		500			{string}	string	"Failed to delete bucket"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets/{bucketName} [delete]
func HandleS3DeleteBucket(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucketName")
		if bucketName == "" {
			http.Error(w, "Bucket name is required", http.StatusBadRequest)
			return
		}

		logger.Info("deleting S3 bucket", "bucket", bucketName)

		_, err := s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			logger.Error("failed to delete S3 bucket", "error", err)
			http.Error(w, fmt.Sprintf("Failed to delete bucket: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success":    true,
			"bucketName": bucketName,
		}

		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleS3ListObjects lists objects in an S3 bucket.
//
//	@Summary		List objects in S3 bucket
//	@Description	Get a list of all objects in an S3 bucket
//	@Tags			aws
//	@Produce		json
//	@Param			bucketName	path		string	true	"Bucket name"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		500			{string}	string	"Failed to list objects"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets/{bucketName}/objects [get]
func HandleS3ListObjects(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucketName")
		if bucketName == "" {
			http.Error(w, "Bucket name is required", http.StatusBadRequest)
			return
		}

		logger.Info("listing objects in S3 bucket", "bucket", bucketName)

		result, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})

		if err != nil {
			logger.Error("failed to list objects", "error", err)
			http.Error(w, fmt.Sprintf("Failed to list objects: %v", err), http.StatusInternalServerError)
			return
		}

		objects := make([]map[string]interface{}, 0, len(result.Contents))
		for _, obj := range result.Contents {
			objects = append(objects, map[string]interface{}{
				"key":          *obj.Key,
				"size":         *obj.Size,
				"lastModified": obj.LastModified,
			})
		}

		response := map[string]interface{}{
			"objects": objects,
			"count":   len(objects),
		}

		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleS3UploadObject uploads an object to S3.
//
//	@Summary		Upload object to S3
//	@Description	Upload a file to an S3 bucket
//	@Tags			aws
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			bucketName	path		string	true	"Bucket name"
//	@Param			file		formData	file	true	"File to upload"
//	@Success		201			{object}	map[string]interface{}
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		500			{string}	string	"Failed to upload file"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets/{bucketName}/objects [post]
func HandleS3UploadObject(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucketName")
		if bucketName == "" {
			http.Error(w, "Bucket name is required", http.StatusBadRequest)
			return
		}

		// Parse multipart form (32MB max)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			logger.Error("failed to parse multipart form", "error", err)
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Error("failed to get file from form", "error", err)
			http.Error(w, "File is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get optional key (filename) from form, default to uploaded filename
		key := r.FormValue("key")
		if key == "" {
			key = header.Filename
		}

		logger.Info("uploading file to S3", "bucket", bucketName, "key", key, "size", header.Size)

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   file,
		})

		if err != nil {
			logger.Error("failed to upload object", "error", err)
			http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"key":     key,
			"bucket":  bucketName,
		}

		if err := encode(w, r, http.StatusCreated, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleS3DeleteObject deletes an object from S3.
//
//	@Summary		Delete object from S3
//	@Description	Delete a file from an S3 bucket
//	@Tags			aws
//	@Produce		json
//	@Param			bucketName	path		string	true	"Bucket name"
//	@Param			key			path		string	true	"Object key"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		500			{string}	string	"Failed to delete object"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets/{bucketName}/objects/{key} [delete]
func HandleS3DeleteObject(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucketName")
		key := r.PathValue("key")

		if bucketName == "" || key == "" {
			http.Error(w, "Bucket name and key are required", http.StatusBadRequest)
			return
		}

		// Decode URL-encoded key
		key = strings.ReplaceAll(key, "%2F", "/")

		logger.Info("deleting object from S3", "bucket", bucketName, "key", key)

		_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})

		if err != nil {
			logger.Error("failed to delete object", "error", err)
			http.Error(w, fmt.Sprintf("Failed to delete object: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"success": true,
			"key":     key,
			"bucket":  bucketName,
		}

		if err := encode(w, r, http.StatusOK, response); err != nil {
			logger.Error("failed to encode response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

// HandleS3GetObject downloads an object from S3.
//
//	@Summary		Download object from S3
//	@Description	Download a file from an S3 bucket
//	@Tags			aws
//	@Produce		octet-stream
//	@Param			bucketName	path		string	true	"Bucket name"
//	@Param			key			path		string	true	"Object key"
//	@Success		200			{file}		binary
//	@Failure		400			{string}	string	"Invalid request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		404			{string}	string	"Object not found"
//	@Failure		500			{string}	string	"Failed to download object"
//	@Security		BearerAuth
//	@Router			/api/v1/aws/s3/buckets/{bucketName}/download/{key} [get]
func HandleS3GetObject(logger *slog.Logger, s3Client *s3.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucketName := r.PathValue("bucketName")
		key := r.PathValue("key")

		if bucketName == "" || key == "" {
			http.Error(w, "Bucket name and key are required", http.StatusBadRequest)
			return
		}

		// Decode URL-encoded key
		key = strings.ReplaceAll(key, "%2F", "/")

		logger.Info("downloading object from S3", "bucket", bucketName, "key", key)

		result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})

		if err != nil {
			logger.Error("failed to get object", "error", err)
			http.Error(w, fmt.Sprintf("Failed to download object: %v", err), http.StatusInternalServerError)
			return
		}
		defer result.Body.Close()

		// Set headers for file download
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", key))
		w.Header().Set("Content-Type", "application/octet-stream")
		if result.ContentLength != nil {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", *result.ContentLength))
		}

		// Stream the file to the response
		_, err = io.Copy(w, result.Body)
		if err != nil {
			logger.Error("failed to stream object", "error", err)
			return
		}
	})
}
