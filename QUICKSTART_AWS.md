# Quick Start: AWS Integration

## Prerequisites

1. AWS Account
2. AWS CLI installed and configured
3. Go 1.21+ installed

## 1. Configure AWS Credentials

### Option A: AWS CLI (Recommended for local development)
```bash
aws configure
```
Enter your:
- AWS Access Key ID
- AWS Secret Access Key
- Default region (e.g., us-east-1)
- Output format (json)

### Option B: Environment Variables
```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your_access_key_here
export AWS_SECRET_ACCESS_KEY=your_secret_key_here
```

### Option C: IAM Role (for EC2/ECS/Lambda)
No configuration needed - the SDK automatically uses the instance role.

## 2. Run the Server

```bash
# Simple run
go run .

# Or with specific region
AWS_REGION=us-west-2 go run .

# Or with profile
AWS_PROFILE=dev go run .
```

## 3. Test AWS Endpoints

### Check if it's working
```bash
# Health check (no AWS required)
curl http://localhost:8080/healthz
```

### List S3 Buckets
```bash
curl http://localhost:8080/api/v1/aws/s3/buckets
```

Expected response:
```json
{
  "buckets": [
    {
      "name": "my-bucket-name",
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

Expected response:
```json
{
  "tables": ["users", "products"],
  "count": 2
}
```

## 4. Create Test Resources (Optional)

### Create an S3 bucket
```bash
aws s3 mb s3://my-test-bucket-$(date +%s)
```

### Create a DynamoDB table
```bash
aws dynamodb create-table \
    --table-name test-table \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1
```

### Test again
```bash
curl http://localhost:8080/api/v1/aws/s3/buckets
curl http://localhost:8080/api/v1/aws/dynamodb/tables
```

## 5. Existing Endpoints (Non-AWS)

These endpoints don't require AWS and work as before:

```bash
# Health check
curl http://localhost:8080/healthz

# Create item
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"This is a test"}'

# Get all items
curl http://localhost:8080/api/v1/items
```

## Troubleshooting

### "Failed to initialize AWS clients"
- Check your AWS credentials are configured
- Verify AWS_REGION is set
- Run `aws sts get-caller-identity` to test credentials

### "Access Denied" errors
- Check IAM permissions for your AWS user/role
- Ensure you have permissions for the service (s3:ListAllMyBuckets, dynamodb:ListTables)
- Example minimal IAM policy:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListAllMyBuckets",
        "dynamodb:ListTables"
      ],
      "Resource": "*"
    }
  ]
}
```

### Server won't start
- Check port 8080 is not already in use
- Try a different port: modify `server.go` line 32

## What Was Added?

### Files Created:
- `aws.go` - AWS client initialization
- `iam_auth.go` - IAM authentication middleware
- `handlers/aws.go` - AWS service handlers

### Files Modified:
- `server.go` - Added AWS client setup
- `routes.go` - Added AWS endpoints
- `go.mod` - Added AWS SDK dependencies

### New Dependencies:
- `github.com/aws/aws-sdk-go-v2` - Core AWS SDK
- `github.com/aws/aws-sdk-go-v2/config` - Configuration
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 client
- `github.com/aws/aws-sdk-go-v2/service/dynamodb` - DynamoDB client

## Next Steps

1. Read the full [AWS_INTEGRATION.md](./AWS_INTEGRATION.md) guide
2. Add more AWS services (Lambda, SQS, SNS, etc.)
3. Enable IAM authentication middleware
4. Deploy to AWS (EC2, ECS, Lambda)
5. Set up CloudWatch monitoring

## Support

For detailed information, see:
- [AWS_INTEGRATION.md](./AWS_INTEGRATION.md) - Complete AWS integration guide
- [RESILIENCE_IMPROVEMENTS.md](./RESILIENCE_IMPROVEMENTS.md) - Server resilience features
- [AWS SDK for Go v2 Docs](https://aws.github.io/aws-sdk-go-v2/)
