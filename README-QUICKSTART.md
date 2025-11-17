# Quick Start Guide

## Prerequisites

Before starting, you need:

1. **Go** (version 1.20+)
2. **Node.js** (version 18+)
3. **AWS Cognito** configured (see below)

## AWS Cognito Setup (REQUIRED)

The server requires AWS Cognito for authentication. You have two options:

### Option 1: Quick Setup with Existing Cognito

If you already have a Cognito User Pool:

```bash
# 1. Copy environment template
cp .env.example .env

# 2. Edit .env and add your Cognito credentials
nano .env  # or your preferred editor

# Required values:
# AWS_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
# AWS_COGNITO_CLIENT_ID=your-client-id
# AWS_COGNITO_CLIENT_SECRET=your-client-secret
```

### Option 2: Create New Cognito User Pool

See detailed guide: **[COGNITO_SETUP.md](./COGNITO_SETUP.md)**

Quick create via AWS CLI:
```bash
# Create User Pool
aws cognito-idp create-user-pool \
  --pool-name "my-app-users" \
  --auto-verified-attributes email \
  --region us-east-1

# Create App Client
aws cognito-idp create-user-pool-client \
  --user-pool-id <YOUR_POOL_ID> \
  --client-name "my-app-client" \
  --generate-secret \
  --explicit-auth-flows ADMIN_NO_SRP_AUTH \
  --region us-east-1

# Copy the IDs to your .env file
```

## Start Development Servers

### One-Command Start (Recommended)

```bash
./start-dev.sh
```

This script will:
- ✅ Check prerequisites (Go, Node.js)
- ✅ Start Go server on http://localhost:8080
- ✅ Start React dev server on http://localhost:5173
- ✅ Set up API proxy automatically
- ✅ Handle cleanup on Ctrl+C

### Manual Start

**Terminal 1 - Go Server:**
```bash
# Make sure .env is configured!
go run cmd/server/main.go
```

**Terminal 2 - React Dev Server:**
```bash
cd web
npm run dev
```

## Test the Setup

1. **Visit the app:** http://localhost:5173

2. **Test the API:**
```bash
curl http://localhost:8080/healthz
# Should return: {"status":"healthy"}
```

3. **Test authentication:**
   - Click "Register" to sign up
   - Check email for verification code
   - Verify email
   - Login with credentials
   - Should see dashboard with user info

## Common Errors & Solutions

### ❌ "AWS_COGNITO_USER_POOL_ID is required"

**Problem:** Missing or empty `.env` file

**Solution:**
```bash
# Copy the template
cp .env.example .env

# Edit and fill in your Cognito credentials
nano .env
```

See [COGNITO_SETUP.md](./COGNITO_SETUP.md) for how to get Cognito credentials.

### ❌ "Maximum update depth exceeded"
✅ **FIXED** - Updated AuthContext to prevent infinite loop

### ❌ "ECONNREFUSED"

**Problem:** Go server is not running

**Solution:**
```bash
# Make sure .env is configured first!
# Then start Go server:
go run cmd/server/main.go

# Or use the start script:
./start-dev.sh
```

### ❌ Port already in use

**Go server (8080):**
```bash
lsof -ti:8080 | xargs kill -9
```

**React dev server (5173):**
```bash
lsof -ti:5173 | xargs kill -9
```

## Environment Configuration

### Go Server (`.env` in project root)

```bash
# AWS Cognito (REQUIRED)
AWS_COGNITO_USER_POOL_ID=us-east-1_XXXXXXXXX
AWS_COGNITO_CLIENT_ID=your-client-id
AWS_COGNITO_CLIENT_SECRET=your-client-secret
AWS_COGNITO_REGION=us-east-1

# AWS General
AWS_REGION=us-east-1
AWS_PROFILE=default

# Server (optional)
SERVER_PORT=8080
SERVER_HOST=localhost
```

### React Frontend (`web/.env`)

Already configured - no changes needed:
```bash
VITE_API_URL=
VITE_APP_NAME=AWS Go Server
```

## Development Workflow

1. **Configure Cognito** (one-time setup)
2. **Start servers:** `./start-dev.sh`
3. **Make changes:** Edit files in `web/src/`
4. **See updates:** React hot-reloads automatically
5. **Test API:** Changes to Go code require server restart

## Build for Production

```bash
# 1. Build React app
cd web
npm run build

# 2. Run Go server (serves React + API)
cd ..
go run cmd/server/main.go

# 3. Visit: http://localhost:8080
```

## Project Structure

```
AWS-Go-Server/
├── .env                    # ⚠️  YOU NEED TO CREATE THIS
├── .env.example            # Template for .env
├── start-dev.sh            # One-command dev start
├── cmd/server/main.go      # Go server entry point
├── internal/
│   ├── config/            # Environment config
│   ├── handlers/          # API handlers
│   ├── auth/             # Cognito integration
│   └── server/           # HTTP server
└── web/                   # React frontend
    ├── src/
    │   ├── api/          # API client
    │   ├── contexts/     # Auth context
    │   ├── pages/        # React pages
    │   └── components/   # React components
    └── dist/             # Built files (production)
```

## Next Steps

1. ✅ **Set up Cognito** - See [COGNITO_SETUP.md](./COGNITO_SETUP.md)
2. ✅ **Start dev servers** - Run `./start-dev.sh`
3. ✅ **Test authentication** - Sign up, verify, login
4. ✅ **Build features** - Add your own pages and API endpoints

## Documentation

- **Cognito Setup:** [COGNITO_SETUP.md](./COGNITO_SETUP.md)
- **React Setup:** [web/SETUP.md](./web/SETUP.md)
- **API Integration:** [web/API_INTEGRATION.md](./web/API_INTEGRATION.md)
- **Architecture:** [web/ARCHITECTURE.md](./web/ARCHITECTURE.md)

## Support

### React Frontend Issues
- Check `web/SETUP.md`
- Check browser console (F12)
- Check Network tab for API errors

### Go Server Issues
- Check server logs in terminal
- Verify `.env` file exists and has correct values
- Test with `curl http://localhost:8080/healthz`

### Cognito Issues
- See [COGNITO_SETUP.md](./COGNITO_SETUP.md)
- Check AWS Cognito console for user pool status
- Verify credentials in `.env` match AWS console
