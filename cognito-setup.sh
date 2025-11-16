#!/bin/bash
# AWS Cognito User Pool Setup Script
# This script creates a Cognito User Pool with app client for the Go server

set -e

# Configuration
POOL_NAME="go-aws-server-users"
CLIENT_NAME="go-aws-server-client"
REGION="${AWS_REGION:-us-east-1}"

echo "Creating Cognito User Pool: $POOL_NAME in region: $REGION"

# Create User Pool
USER_POOL_OUTPUT=$(aws cognito-idp create-user-pool \
    --pool-name "$POOL_NAME" \
    --region "$REGION" \
    --policies '{
        "PasswordPolicy": {
            "MinimumLength": 8,
            "RequireUppercase": true,
            "RequireLowercase": true,
            "RequireNumbers": true,
            "RequireSymbols": false
        }
    }' \
    --auto-verified-attributes email \
    --username-attributes email \
    --mfa-configuration OFF \
    --account-recovery-setting '{
        "RecoveryMechanisms": [
            {
                "Priority": 1,
                "Name": "verified_email"
            }
        ]
    }' \
    --schema '[
        {
            "Name": "email",
            "AttributeDataType": "String",
            "Required": true,
            "Mutable": true
        },
        {
            "Name": "name",
            "AttributeDataType": "String",
            "Required": false,
            "Mutable": true
        }
    ]' \
    --user-attribute-update-settings '{
        "AttributesRequireVerificationBeforeUpdate": ["email"]
    }' \
    --output json)

USER_POOL_ID=$(echo "$USER_POOL_OUTPUT" | jq -r '.UserPool.Id')
echo "✓ User Pool created: $USER_POOL_ID"

# Create User Pool Client
CLIENT_OUTPUT=$(aws cognito-idp create-user-pool-client \
    --user-pool-id "$USER_POOL_ID" \
    --client-name "$CLIENT_NAME" \
    --region "$REGION" \
    --generate-secret \
    --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH ALLOW_USER_SRP_AUTH \
    --token-validity-units '{
        "AccessToken": "minutes",
        "IdToken": "minutes",
        "RefreshToken": "days"
    }' \
    --access-token-validity 60 \
    --id-token-validity 60 \
    --refresh-token-validity 30 \
    --prevent-user-existence-errors ENABLED \
    --output json)

CLIENT_ID=$(echo "$CLIENT_OUTPUT" | jq -r '.UserPoolClient.ClientId')
CLIENT_SECRET=$(echo "$CLIENT_OUTPUT" | jq -r '.UserPoolClient.ClientSecret')
echo "✓ User Pool Client created: $CLIENT_ID"

# Get User Pool Domain prefix (for issuer URL)
ISSUER_URL="https://cognito-idp.$REGION.amazonaws.com/$USER_POOL_ID"

echo ""
echo "========================================="
echo "AWS Cognito Setup Complete!"
echo "========================================="
echo ""
echo "Add these to your .env file:"
echo ""
echo "AWS_COGNITO_REGION=$REGION"
echo "AWS_COGNITO_USER_POOL_ID=$USER_POOL_ID"
echo "AWS_COGNITO_CLIENT_ID=$CLIENT_ID"
echo "AWS_COGNITO_CLIENT_SECRET=$CLIENT_SECRET"
echo ""
echo "JWT Token Issuer URL:"
echo "$ISSUER_URL"
echo ""
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Add the above environment variables to your .env file"
echo "2. The Go application will validate JWT tokens from this User Pool"
echo "3. Users can sign up and login through the API endpoints"
echo ""
