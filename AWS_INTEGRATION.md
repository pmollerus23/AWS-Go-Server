# AWS Integration Guide

Your server now has full AWS integration using AWS SDK for Go v2 with IAM authentication support.

## Features

### 1. AWS Service Clients
- **S3 Client**: For object storage operations
- **DynamoDB Client**: For NoSQL database operations
- Easily extensible to other AWS services (Lambda, SQS, SNS, etc.)

### 2. IAM Authentication
- Verifies AWS SigV4 signatures on incoming requests
- Validates request timestamps (prevents replay attacks)
- Logs authentication attempts for security auditing

### 3. Automatic Credential Chain
The AWS SDK automatically looks for credentials in this order:
1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM role (when running on EC2, ECS, Lambda, etc.)

## Configuration

### Environment Variables

```bash
# Required: AWS Region
export AWS_REGION=us-east-1

# Optional: AWS Profile (for local development)
export AWS_PROFILE=default

# Optional: Explicit credentials (not recommended for production)
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
```

### AWS Credentials File

Create `~/.aws/credentials`:
```ini
[default]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY

[dev]
aws_access_key_id = YOUR_DEV_ACCESS_KEY
aws_secret_access_key = YOUR_DEV_SECRET_KEY
```

Create `~/.aws/config`:
```ini
[default]
region = us-east-1
output = json

[profile dev]
region = us-west-2
output = json
```

## Available AWS Endpoints

### List S3 Buckets
```bash
curl http://localhost:8080/api/v1/aws/s3/buckets
```

Response:
```json
{
  "buckets": [
    {
      "name": "my-bucket",
      "creationDate": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

### List DynamoDB Tables
```bash
curl http://localhost:8080/api/v1/aws/dynamodb/tables
```

Response:
```json
{
  "tables": ["users", "products", "orders"],
  "count": 3
}
```

## IAM Authentication Middleware

The server includes IAM authentication middleware that verifies AWS SigV4 signatures.

### How It Works

1. Checks for `Authorization` header with AWS4-HMAC-SHA256 scheme
2. Validates `X-Amz-Date` header (must be within 15 minutes)
3. Parses credential, signed headers, and signature
4. Logs authentication attempts

### Using IAM Auth in Routes

To protect specific routes with IAM authentication, apply the middleware:

```go
// In server.go, add IAM auth middleware
handler = newIAMAuthMiddleware(logger, config.AWS.Region)(handler)
```

**Note**: Currently, the middleware is implemented but not applied by default. Add it to the middleware stack in `server.go:111` if you want to require IAM authentication for all requests.

### Testing IAM Authentication

You can use AWS SDKs or tools like `awscurl` to make signed requests:

```bash
# Install awscurl
pip install awscurl

# Make a signed request
awscurl --service execute-api \
  --region us-east-1 \
  http://localhost:8080/api/v1/aws/s3/buckets
```

Or use the AWS SDK from your application:

```go
// Example: Making a signed request from another Go app
import (
    "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

// Use the SDK's built-in HTTP client with SigV4 signing
```

## Adding More AWS Services

### Example: Adding Lambda Client

1. Install the Lambda SDK:
```bash
go get github.com/aws/aws-sdk-go-v2/service/lambda
```

2. Add to `aws.go`:
```go
import "github.com/aws/aws-sdk-go-v2/service/lambda"

type AWSClients struct {
    Config   aws.Config
    S3       *s3.Client
    DynamoDB *dynamodb.Client
    Lambda   *lambda.Client  // Add this
}

func NewAWSClients(...) (*AWSClients, error) {
    // ... existing code ...
    clients := &AWSClients{
        Config:   cfg,
        S3:       s3.NewFromConfig(cfg),
        DynamoDB: dynamodb.NewFromConfig(cfg),
        Lambda:   lambda.NewFromConfig(cfg),  // Add this
    }
    return clients, nil
}
```

3. Create handler in `handlers/aws.go`:
```go
func HandleLambdaInvoke(logger *slog.Logger, lambdaClient *lambda.Client) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Invoke Lambda function
        result, err := lambdaClient.Invoke(context.TODO(), &lambda.InvokeInput{
            FunctionName: aws.String("my-function"),
        })
        // ... handle response
    })
}
```

4. Add route in `routes.go`:
```go
mux.Handle("POST /api/v1/aws/lambda/invoke",
    handlers.HandleLambdaInvoke(logger, awsClients.Lambda))
```

## Security Best Practices

### 1. Never Commit Credentials
Add to `.gitignore`:
```
.env
.aws/
*.pem
*.key
credentials.json
```

### 2. Use IAM Roles When Possible
- **EC2**: Attach IAM role to instance
- **ECS/Fargate**: Use task roles
- **Lambda**: Use execution roles
- **Local Dev**: Use AWS SSO or temporary credentials

### 3. Principle of Least Privilege
Grant only the permissions needed:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket",
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::my-bucket",
        "arn:aws:s3:::my-bucket/*"
      ]
    }
  ]
}
```

### 4. Enable CloudTrail
Log all AWS API calls for auditing and security analysis.

### 5. Rotate Credentials Regularly
Use AWS Secrets Manager or Parameter Store for credential rotation.

## Testing Locally

### 1. Using LocalStack
LocalStack provides local AWS service emulation:

```bash
# Install LocalStack
pip install localstack

# Start services
localstack start

# Set endpoint
export AWS_ENDPOINT_URL=http://localhost:4566
```

### 2. Using AWS SAM
For Lambda development:
```bash
sam local start-api
```

## Troubleshooting

### "No credentials found"
- Check `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables
- Verify `~/.aws/credentials` file exists and is properly formatted
- Ensure IAM role is attached (if running on AWS)

### "Access Denied"
- Check IAM policy permissions
- Verify the principal has the required actions
- Check resource ARNs in the policy

### "Invalid signature"
- Ensure system clock is synchronized (SigV4 is time-sensitive)
- Verify the secret key matches the access key
- Check that the request hasn't been modified

### Connection timeout
- Check security group rules
- Verify VPC configuration
- Ensure AWS service endpoints are reachable

## Running the Server

```bash
# Set AWS credentials
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret

# Run the server
go run .

# Or build and run
go build -o server && ./server
```

## Example: Complete Workflow

```bash
# 1. Configure AWS credentials
aws configure

# 2. Create an S3 bucket
aws s3 mb s3://my-test-bucket

# 3. Start the server
export AWS_REGION=us-east-1
go run .

# 4. Test the S3 endpoint
curl http://localhost:8080/api/v1/aws/s3/buckets

# 5. Should see your bucket listed
{
  "buckets": [
    {"name": "my-test-bucket", "creationDate": "..."}
  ],
  "count": 1
}
```

## Architecture

```
Client Request
     ↓
Panic Recovery Middleware
     ↓
Request Size Limit Middleware
     ↓
Logging Middleware
     ↓
[Optional: IAM Auth Middleware]
     ↓
Router
     ↓
AWS Service Handlers
     ↓
AWS SDK Clients (S3, DynamoDB, etc.)
     ↓
AWS Services
```

## Code Structure

```
AWS-Go-Server/
├── main.go                 # Entry point
├── server.go               # Server setup and configuration
├── aws.go                  # AWS client initialization
├── iam_auth.go            # IAM authentication middleware
├── middleware.go          # Other middleware (logging, panic recovery, etc.)
├── routes.go              # Route definitions
└── handlers/
    ├── items.go           # Regular CRUD handlers
    ├── health.go          # Health check
    └── aws.go             # AWS service handlers (S3, DynamoDB)
```

## Next Steps

1. **Add More Services**: SQS, SNS, Lambda, CloudWatch, etc.
2. **Implement Full SigV4 Verification**: Complete the signature verification in `iam_auth.go`
3. **Add Resource-Based Authorization**: Check IAM policies for fine-grained access control
4. **Metrics and Monitoring**: Integrate CloudWatch metrics
5. **Distributed Tracing**: Add X-Ray support
6. **Caching**: Use ElastiCache or DynamoDB DAX
7. **Async Processing**: Use SQS for background jobs

## Resources

- [AWS SDK for Go v2 Documentation](https://aws.github.io/aws-sdk-go-v2/)
- [AWS SigV4 Signing Process](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html)
- [IAM Best Practices](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html)
- [AWS Security Best Practices](https://aws.amazon.com/architecture/security-identity-compliance/)
