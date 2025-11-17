# API Integration Guide

This document explains how the React frontend is integrated with the Go server backend.

## API Endpoints

The React app connects to the following Go server endpoints:

### Authentication Endpoints

| Endpoint | Method | Description | Request | Response |
|----------|--------|-------------|---------|----------|
| `/api/v1/auth/signup` | POST | Register new user | `{ email, password, name? }` | `{ message, email }` |
| `/api/v1/auth/confirm` | POST | Confirm email | `{ email, code }` | `{ message }` |
| `/api/v1/auth/login` | POST | Login | `{ email, password }` | `{ message, tokens }` |
| `/api/v1/auth/refresh` | POST | Refresh tokens | `{ refresh_token, email }` | `{ message, tokens }` |
| `/api/v1/auth/forgot-password` | POST | Request password reset | `{ email }` | `{ message }` |
| `/api/v1/auth/reset-password` | POST | Reset password | `{ email, code, new_password }` | `{ message }` |

### Protected Endpoints (Require Authentication)

| Endpoint | Method | Description | Request | Response |
|----------|--------|-------------|---------|----------|
| `/api/v1/items` | GET | Get all items | - | `Item[]` |
| `/api/v1/items` | POST | Create item | `{ name, description }` | `{ id, name, description }` |

## Authentication Flow

### 1. Sign Up Flow

```typescript
// User signs up
await authApi.signUp({ email, password, name });

// User receives email with verification code
// User confirms email
await authApi.confirmSignUp({ email, code });

// User can now log in
```

### 2. Login Flow

```typescript
// User logs in
const response = await authApi.login({ email, password });

// Response contains Cognito tokens:
{
  message: "Login successful",
  tokens: {
    access_token: "...",  // For API authorization
    id_token: "...",      // Contains user info
    refresh_token: "...", // For refreshing tokens
    expires_in: 3600,
    token_type: "Bearer"
  }
}

// Tokens are stored in localStorage
// User info is extracted from ID token JWT
```

### 3. Token Management

The app automatically:
- Stores tokens in localStorage
- Extracts user info from ID token
- Injects access token in API requests
- Checks token expiration on app load
- Refreshes tokens when expired
- Clears tokens on logout

### 4. API Request Flow

```typescript
// Protected API call
const items = await itemsApi.getAll();

// Under the hood:
1. API client retrieves access_token from localStorage
2. Adds Authorization header: "Bearer <access_token>"
3. Makes request to /api/v1/items
4. Go server validates JWT token
5. Returns data if authorized
```

## Token Storage

Tokens are stored in localStorage with the following keys:

```typescript
{
  "access_token": "eyJhbG...",
  "id_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "user_email": "user@example.com"
}
```

## User Information Extraction

User information is extracted from the ID token JWT:

```typescript
// ID token contains Cognito claims:
{
  sub: "user-uuid",               // User ID
  email: "user@example.com",
  email_verified: true,
  name: "John Doe",
  given_name: "John",
  family_name: "Doe",
  "cognito:username": "johndoe",
  "cognito:groups": ["admin"],
  exp: 1234567890,
  iat: 1234567890
}

// Mapped to User object:
{
  id: "user-uuid",
  email: "user@example.com",
  username: "johndoe",
  name: "John Doe",
  firstName: "John",
  lastName: "Doe",
  emailVerified: true,
  groups: ["admin"]
}
```

## Using the API in Components

### useQuery Hook (GET requests)

```typescript
import { useQuery } from '../hooks';
import { itemsApi } from '../api';

const { data, isLoading, isError, error, refetch } = useQuery(
  async () => itemsApi.getAll(),
  {
    enabled: true,              // Only run if true
    refetchOnMount: true,       // Refetch when component mounts
    refetchInterval: 30000,     // Poll every 30 seconds
    onSuccess: (data) => { },   // Callback on success
    onError: (error) => { },    // Callback on error
    retry: 3,                   // Retry failed requests
    retryDelay: 1000,           // Wait 1s between retries
  }
);
```

### useMutation Hook (POST/PUT/DELETE requests)

```typescript
import { useMutation } from '../hooks';
import { itemsApi } from '../api';

const mutation = useMutation(
  async (data) => itemsApi.create(data),
  {
    onSuccess: (result) => {
      console.log('Item created:', result);
    },
    onError: (error) => {
      console.error('Failed:', error);
    },
  }
);

// Trigger mutation
mutation.mutate({ name: 'New Item', description: 'Description' });

// Or await result
const result = await mutation.mutateAsync({ name: 'Item', description: 'Desc' });
```

### Auth Context

```typescript
import { useAuth } from '../contexts';

const { user, isAuthenticated, isLoading, login, logout } = useAuth();

// Login
await login({ email, password });

// Logout
await logout();

// Check auth status
if (isAuthenticated) {
  console.log('User:', user);
}
```

## Error Handling

API errors are typed and contain:

```typescript
interface ApiError {
  message: string;
  status?: number;    // HTTP status code
  code?: string;      // Error code
  details?: object;   // Additional error details
}
```

Example:

```typescript
try {
  await authApi.login({ email, password });
} catch (error: any) {
  if (error.status === 401) {
    console.log('Invalid credentials');
  } else if (error.message) {
    console.log('Error:', error.message);
  }
}
```

## Environment Configuration

Set the API URL in `.env`:

```bash
# Use relative URLs (same server)
VITE_API_URL=

# Or specify full URL
VITE_API_URL=http://localhost:8080
```

## Development

1. Start Go server: `cd .. && go run cmd/server/main.go`
2. Start React dev server: `npm run dev`
3. Access app at `http://localhost:5173`

## Production

1. Build React app: `npm run build`
2. Go server serves from `web/dist/`
3. All API routes prefixed with `/api/v1/`
4. SPA routing handled by Go server fallback

## Security Notes

- Access tokens are JWT validated by Go server
- Tokens stored in localStorage (consider httpOnly cookies for production)
- CORS handled by Go server
- All protected routes require valid JWT
- Tokens auto-refresh before expiration
- User cannot modify JWT claims (server-side validation)
