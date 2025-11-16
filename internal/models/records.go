package models

// DynamoDBRecord represents a record to be stored in DynamoDB.
type DynamoDBRecord struct {
	ID        int    `json:"id" dynamodbav:"id" example:"1"`
	Name      string `json:"name" dynamodbav:"name" example:"Sample Record"`
	UpdatedAt int64  `json:"updated_at" dynamodbav:"updated_at" example:"1699999999"`
}
