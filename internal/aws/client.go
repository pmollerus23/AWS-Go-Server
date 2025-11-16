package aws

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appConfig "github.com/pmollerus23/go-aws-server/internal/config"
)

// Clients holds all AWS service clients.
type Clients struct {
	Config   aws.Config
	S3       *s3.Client
	DynamoDB *dynamodb.Client
	Cognito  *cognito.Client
}

// NewClients creates and initializes AWS service clients.
func NewClients(ctx context.Context, logger *slog.Logger, awsConfig appConfig.AWSConfig) (*Clients, error) {
	// Load AWS configuration
	// This will use the default credential chain:
	// 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)
	// 2. Shared credentials file (~/.aws/credentials)
	// 3. IAM role (if running on EC2, ECS, Lambda, etc.)
	var configOpts []func(*config.LoadOptions) error

	if awsConfig.Region != "" {
		configOpts = append(configOpts, config.WithRegion(awsConfig.Region))
	}

	if awsConfig.Profile != "" {
		configOpts = append(configOpts, config.WithSharedConfigProfile(awsConfig.Profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		logger.Error("failed to load AWS config", "error", err)
		return nil, err
	}

	logger.Info("AWS config loaded",
		"region", cfg.Region,
	)

	// Create service clients
	clients := &Clients{
		Config:   cfg,
		S3:       s3.NewFromConfig(cfg),
		DynamoDB: dynamodb.NewFromConfig(cfg),
		Cognito:  cognito.NewFromConfig(cfg),
	}

	return clients, nil
}
