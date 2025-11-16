# AWS Cognito Setup Guide

This guide will help you set up AWS Cognito for authentication in your Go server.

## Prerequisites

- AWS CLI installed and configured with appropriate credentials
- AWS account with permissions to create Cognito resources
- `jq` installed (for the bash script method)

## Setup Options

### Option 1: Bash Script (Quick Setup)

```bash
# Make the script executable
chmod +x cognito-setup.sh

# Run the script (optionally set AWS_REGION)
AWS_REGION=us-east-1 ./cognito-setup.sh
```

The script will output environment variables that you need to add to your `.env` file.

### Option 2: CloudFormation (Infrastructure as Code)

```bash
# Deploy the CloudFormation stack
aws cloudformation create-stack \
  --stack-name go-aws-server-cognito \
  --template-body file://cognito-cloudformation.yaml \
  --region us-east-1

# Wait for stack creation to complete
aws cloudformation wait stack-create-complete \
  --stack-name go-aws-server-cognito \
  --region us-east-1

# Get the outputs
aws cloudformation describe-stacks \
  --stack-name go-aws-server-cognito \
  --region us-east-1 \
  --query 'Stacks[0].Outputs'
```

After deployment, get the Client Secret:

```bash
# First, get the User Pool ID and Client ID from stack outputs
USER_POOL_ID=$(aws cloudformation describe-stacks \
  --stack-name go-aws-server-cognito \
  --query 'Stacks[0].Outputs[?OutputKey==`UserPoolId`].OutputValue' \
  --output text)

CLIENT_ID=$(aws cloudformation describe-stacks \
  --stack-name go-aws-server-cognito \
  --query 'Stacks[0].Outputs[?OutputKey==`UserPoolClientId`].OutputValue' \
  --output text)

# Get the Client Secret
CLIENT_SECRET=$(aws cognito-idp describe-user-pool-client \
  --user-pool-id "$USER_POOL_ID" \
  --client-id "$CLIENT_ID" \
  --query 'UserPoolClient.ClientSecret' \
  --output text)

echo "AWS_COGNITO_CLIENT_SECRET=$CLIENT_SECRET"
```

## Environment Configuration

Create a `.env` file in your project root with the following variables:

```bash
# AWS Configuration
AWS_REGION=us-east-1

# Cognito Configuration
AWS_COGNITO_REGION=us-east-1
AWS_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
AWS_COGNITO_CLIENT_ID=your-client-id-here
AWS_COGNITO_CLIENT_SECRET=your-client-secret-here
```

## What Gets Created

The setup creates:

1. **Cognito User Pool** with:
   - Email-based authentication
   - Auto-verified email addresses
   - Password policy (min 8 chars, uppercase, lowercase, numbers)
   - Email recovery mechanism

2. **User Pool Client** with:
   - Client secret for server-side authentication
   - Password-based authentication flow
   - Refresh token support
   - Access token validity: 60 minutes
   - ID token validity: 60 minutes
   - Refresh token validity: 30 days

## Testing Your Setup

After setting up Cognito and configuring your application:

1. **Sign up a test user:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test1234",
    "name": "Test User"
  }'
```

2. **Confirm the user** (check email for verification code or use AWS CLI):
```bash
aws cognito-idp admin-confirm-sign-up \
  --user-pool-id <USER_POOL_ID> \
  --username test@example.com
```

3. **Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test1234"
  }'
```

4. **Use the access token:**
```bash
curl -X GET http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer <access_token>"
```

## Cleanup

To delete the Cognito resources:

### If using bash script:
```bash
# Delete user pool (this also deletes the client)
aws cognito-idp delete-user-pool --user-pool-id <USER_POOL_ID>
```

### If using CloudFormation:
```bash
aws cloudformation delete-stack --stack-name go-aws-server-cognito
```

## Next Steps

1. Run the setup script or deploy the CloudFormation template
2. Add the output values to your `.env` file
3. Restart your Go application
4. Test the authentication endpoints
