# AWS Cognito Authentication Integration

This document describes the AWS Cognito authentication integration in your Go server.

## Overview

Your Go server now uses AWS Cognito for authentication and authorization. All `/api/v1/*` endpoints (except auth endpoints) are protected and require a valid JWT access token from Cognito.

## What Was Added

### 1. AWS Cognito Setup Scripts
- `cognito-setup.sh` - Bash script to create Cognito User Pool via AWS CLI
- `cognito-cloudformation.yaml` - CloudFormation template for infrastructure as code
- `COGNITO_SETUP.md` - Detailed setup instructions

### 2. Configuration
- **Config Updated**: `internal/config/config.go` now includes Cognito settings
- **Environment Variables**: Added `.env.example` with required Cognito configuration

### 3. Authentication Service
- **Cognito Service**: `internal/auth/cognito.go` - Handles all Cognito operations:
  - User signup
  - Email verification
  - Login/authentication
  - Token refresh
  - Password reset
  - JWT token validation using JWKS

### 4. Authentication Middleware
- **Auth Middleware**: `internal/middleware/auth.go` - Validates JWT tokens on protected routes
- **Permission Middleware**: Check user permissions and roles
- **Admin Middleware**: Restrict access to admin users only

### 5. Authentication Handlers
- **Auth Handlers**: `internal/handlers/auth.go` - HTTP handlers for:
  - `POST /api/v1/auth/signup` - Register new user
  - `POST /api/v1/auth/confirm` - Verify email with code
  - `POST /api/v1/auth/login` - Authenticate and get tokens
  - `POST /api/v1/auth/refresh` - Refresh access token
  - `POST /api/v1/auth/forgot-password` - Request password reset
  - `POST /api/v1/auth/reset-password` - Confirm password reset

### 6. Protected Routes
All existing API endpoints are now protected:
- `GET /api/v1/items` - Requires authentication
- `POST /api/v1/items` - Requires authentication
- `GET /api/v1/aws/s3/buckets` - Requires authentication
- `GET /api/v1/aws/dynamodb/tables` - Requires authentication

### 7. Swagger Documentation
- Updated Swagger annotations with Bearer authentication
- All auth endpoints documented with request/response schemas
- Protected endpoints marked with `@Security BearerAuth`

## Quick Start

### 1. Set Up AWS Cognito User Pool

Choose one of the following methods:

**Option A: Using the bash script**
```bash
chmod +x cognito-setup.sh
AWS_REGION=us-east-1 ./cognito-setup.sh
```

**Option B: Using CloudFormation**
```bash
aws cloudformation create-stack \
  --stack-name go-aws-server-cognito \
  --template-body file://cognito-cloudformation.yaml \
  --region us-east-1
```

See `COGNITO_SETUP.md` for detailed instructions.

### 2. Configure Environment Variables

Copy the values from the setup script output to your `.env` file:

```bash
# AWS Configuration
AWS_REGION=us-east-1

# Cognito Configuration
AWS_COGNITO_REGION=us-east-1
AWS_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
AWS_COGNITO_CLIENT_ID=your-client-id
AWS_COGNITO_CLIENT_SECRET=your-client-secret
```

### 3. Run the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Usage Examples

### 1. Sign Up a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123",
    "name": "John Doe"
  }'
```

**Response:**
```json
{
  "message": "User registered successfully. Please check your email for verification code.",
  "email": "user@example.com"
}
```

### 2. Confirm Email Verification

Check your email for the verification code, then:

```bash
curl -X POST http://localhost:8080/api/v1/auth/confirm \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456"
  }'
```

**Alternative:** Confirm via AWS CLI (for testing):
```bash
aws cognito-idp admin-confirm-sign-up \
  --user-pool-id <USER_POOL_ID> \
  --username user@example.com
```

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

**Response:**
```json
{
  "message": "Login successful",
  "tokens": {
    "access_token": "eyJraWQiOiJ...",
    "id_token": "eyJraWQiOiJ...",
    "refresh_token": "eyJjdHkiOi...",
    "expires_in": 3600,
    "token_type": "Bearer"
  }
}
```

### 4. Access Protected Endpoints

Use the `access_token` from the login response:

```bash
curl -X GET http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer eyJraWQiOiJ..."
```

### 5. Create an Item (Protected)

```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer eyJraWQiOiJ..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Item",
    "description": "Item description"
  }'
```

### 6. Refresh Access Token

When your access token expires (after 60 minutes):

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJjdHkiOi...",
    "email": "user@example.com"
  }'
```

### 7. Password Reset

Request password reset:
```bash
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

Confirm password reset with code from email:
```bash
curl -X POST http://localhost:8080/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "code": "123456",
    "new_password": "NewSecurePass123"
  }'
```

## Swagger UI

Access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

You can test all endpoints directly from the Swagger UI. Use the "Authorize" button to add your Bearer token.

## Architecture

### Authentication Flow

1. **User Registration**
   - User submits email/password to `/api/v1/auth/signup`
   - Cognito creates user and sends verification email
   - User confirms email with code via `/api/v1/auth/confirm`

2. **Login**
   - User submits credentials to `/api/v1/auth/login`
   - Server calls Cognito's `InitiateAuth` API
   - Cognito returns access token, ID token, and refresh token
   - Client stores access token for subsequent requests

3. **Protected Route Access**
   - Client includes `Authorization: Bearer <access_token>` header
   - Auth middleware extracts and validates token
   - Token validated against Cognito's JWKS (public keys)
   - User information extracted from token claims
   - User object added to request context
   - Handler can access user via `auth.GetUser(ctx)`

4. **Token Refresh**
   - When access token expires, client uses refresh token
   - Server calls Cognito's refresh token API
   - New access and ID tokens returned

### JWT Token Validation

The server validates JWT tokens using AWS Cognito's JWKS (JSON Web Key Set):

1. Token extracted from `Authorization` header
2. Token signature verified using public keys from Cognito
3. Token expiration checked
4. Token issuer verified (must be your Cognito User Pool)
5. Claims extracted (user ID, email, username, roles)
6. User object created and added to request context

### Security Features

- **HMAC Secret Hash**: All Cognito API calls use secret hash for enhanced security
- **JWT Signature Verification**: Tokens validated using Cognito's public keys
- **JWKS Caching**: Public keys cached for 1 hour to reduce latency
- **Token Expiration**: Access tokens expire after 60 minutes
- **Refresh Tokens**: Long-lived (30 days) for seamless token renewal
- **Role-Based Access Control**: Support for user roles and permissions
- **Password Policy**: Enforced by Cognito (min 8 chars, uppercase, lowercase, numbers)

## Roles and Permissions

The system supports role-based access control:

### Predefined Roles
- **user** - Read-only access to items
- **editor** - Read and write access to items
- **admin** - Full access to all resources

### Using Roles

To require specific permissions on a route:
```go
mux.Handle("DELETE /api/v1/items/{id}",
  authMiddleware(
    middleware.RequirePermission(auth.PermissionDeleteItems, logger)(
      handlers.HandleItemsDelete(logger),
    ),
  ),
)
```

To require admin access:
```go
mux.Handle("GET /api/v1/admin/users",
  authMiddleware(
    middleware.RequireAdmin(logger)(
      handlers.HandleAdminListUsers(logger),
    ),
  ),
)
```

### Assigning Roles to Users

Roles are assigned via Cognito Groups:

```bash
# Create a group
aws cognito-idp create-group \
  --group-name admin \
  --user-pool-id <USER_POOL_ID> \
  --description "Administrator group"

# Add user to group
aws cognito-idp admin-add-user-to-group \
  --user-pool-id <USER_POOL_ID> \
  --username user@example.com \
  --group-name admin
```

Groups appear in the JWT token as `cognito:groups` claim.

## Troubleshooting

### Common Issues

1. **"Unauthorized: missing authorization header"**
   - Make sure you're including the `Authorization` header
   - Format: `Authorization: Bearer <access_token>`

2. **"Unauthorized: invalid token"**
   - Token may be expired (access tokens last 60 minutes)
   - Use refresh token to get a new access token
   - Verify you're using the `access_token`, not `id_token`

3. **"User email not verified"**
   - Confirm email with verification code
   - Or use AWS CLI: `aws cognito-idp admin-confirm-sign-up`

4. **Environment variable errors on startup**
   - Ensure all required Cognito env vars are set in `.env`
   - Check that values match your Cognito User Pool

5. **"Failed to refresh JWKS"**
   - Network issue connecting to Cognito
   - Check AWS region matches your User Pool region
   - Verify internet connectivity

## Next Steps

### Recommended Enhancements

1. **Add User Management Endpoints**
   - Get current user profile
   - Update user profile
   - Change password

2. **Implement Database Storage**
   - Store user metadata in DynamoDB
   - Link items to user IDs
   - Add multi-tenancy support

3. **Add More Granular Permissions**
   - Item-level permissions
   - Resource ownership checks
   - Custom permission logic

4. **Implement Rate Limiting**
   - Protect against brute force attacks
   - Limit API calls per user

5. **Add Monitoring**
   - Log authentication events
   - Track failed login attempts
   - Monitor token validation failures

## Resources

- [AWS Cognito Documentation](https://docs.aws.amazon.com/cognito/)
- [Cognito User Pools](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-identity-pools.html)
- [JWT Specification](https://jwt.io/)
- [Swagger Documentation](http://localhost:8080/swagger/index.html)

## Support

For issues or questions:
1. Check the `COGNITO_SETUP.md` for setup instructions
2. Review error logs in console output
3. Verify environment variables are correct
4. Test authentication flow using Swagger UI
