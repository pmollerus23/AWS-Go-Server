# React Frontend Setup Guide

## Quick Start

### Development Mode (with Hot Reload)

**Terminal 1 - Start Go Server:**
```bash
cd /home/flip/Projects/Go_Learning/AWS-Go-Server
go run cmd/server/main.go
```

**Terminal 2 - Start React Dev Server:**
```bash
cd web
npm run dev
```

Then visit **http://localhost:5173**

The Vite dev server will proxy API requests to the Go server at `localhost:8080`.

### Production Mode (Single Server)

**Build React App:**
```bash
cd web
npm run build
```

**Start Go Server:**
```bash
cd ..
go run cmd/server/main.go
```

Then visit **http://localhost:8080**

The Go server serves both the React app and the API.

## How It Works

### API Client Configuration

The API client automatically handles both development and production:

**Empty baseURL (default):**
- Uses relative URLs: `/api/v1/auth/login`
- Works with Vite proxy in dev
- Works with Go server in production

**Custom baseURL (optional):**
```bash
# .env
VITE_API_URL=http://localhost:8080
```
- Uses absolute URLs: `http://localhost:8080/api/v1/auth/login`
- Useful for external API servers
- Requires CORS configuration

### Vite Proxy (Development)

Configured in `vite.config.ts`:
```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
    '/healthz': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

This proxies all `/api/*` requests from `localhost:5173` to `localhost:8080`.

### Go Server (Production)

The Go server at `internal/server/routes.go` serves:
1. **API routes:** `/api/v1/*` - JSON API endpoints
2. **Static files:** `/*` - React app from `web/dist/`
3. **SPA fallback:** Serves `index.html` for unknown routes (client-side routing)

## Testing Authentication Flow

### 1. Sign Up (if AWS Cognito is configured)

Visit: http://localhost:5173/register

- Enter email and password
- Submit form
- Check email for verification code
- Enter code to confirm
- Redirected to login

### 2. Login

Visit: http://localhost:5173/login

- Enter credentials
- Submit form
- On success: Redirected to home page
- Tokens stored in localStorage
- User info displayed

### 3. Authenticated Requests

Home page automatically:
- Fetches items from `/api/v1/items`
- Includes `Authorization: Bearer <token>` header
- Displays items or shows login prompt

### 4. Logout

Click "Logout" in header:
- Clears tokens from localStorage
- Updates auth state
- Redirects to home (shows login prompt)

## Troubleshooting

### API Calls Failing

**Check Go server is running:**
```bash
curl http://localhost:8080/healthz
# Should return: {"status":"healthy"}
```

**Check browser console:**
- Open DevTools → Network tab
- Try logging in
- Look for `/api/v1/auth/login` request
- Check request URL, headers, and response

**Common issues:**
- Go server not running → Start it first
- Wrong port → Check Go server logs for actual port
- CORS errors → Not needed with proxy/same-origin

### Routing Not Working

**Symptoms:**
- Clicking links causes page refresh
- 404 errors on direct URL access

**Solutions:**
- Development: Vite handles routing automatically
- Production: Go server's SPA handler serves `index.html` for all routes

**Test routing:**
```bash
# Should return index.html for all these:
curl http://localhost:8080/
curl http://localhost:8080/login
curl http://localhost:8080/profile
curl http://localhost:8080/any-route

# Should return JSON:
curl http://localhost:8080/api/v1/items
curl http://localhost:8080/healthz
```

### Token Issues

**Check tokens in localStorage:**
- Open DevTools → Application → Local Storage → http://localhost:5173
- Should see: `access_token`, `id_token`, `refresh_token`, `user_email`

**Clear tokens manually:**
```javascript
// In browser console:
localStorage.clear()
location.reload()
```

**Token expiration:**
- Access tokens expire (usually 1 hour)
- App automatically tries to refresh
- If refresh fails, redirects to login

## File Structure

```
web/
├── src/
│   ├── api/              # API client and endpoints
│   │   ├── client.ts     # ✅ Fixed buildURL for relative URLs
│   │   ├── auth.api.ts   # Auth endpoints
│   │   └── items.api.ts  # Items endpoints
│   ├── contexts/
│   │   └── AuthContext.tsx  # Token management & auth state
│   ├── hooks/
│   │   ├── useQuery.ts      # GET requests
│   │   └── useMutation.ts   # POST/PUT/DELETE requests
│   ├── pages/
│   │   ├── HomePage.tsx     # Dashboard with items
│   │   ├── LoginPage.tsx    # Login form
│   │   └── SignUpPage.tsx   # Registration + verification
│   └── App.tsx              # ✅ Routes with React Router
├── .env                  # Environment variables (empty baseURL)
├── vite.config.ts        # ✅ Proxy configuration
└── package.json

../internal/server/
├── routes.go             # ✅ API routes + SPA handler
└── server.go             # Server setup
```

## Environment Variables

**Development (.env):**
```bash
# Empty = use relative URLs (recommended)
VITE_API_URL=

# App name
VITE_APP_NAME=AWS Go Server
```

**Production:**
No environment variables needed - uses same origin.

## Build Output

```bash
npm run build

# Output:
dist/
├── index.html              # 0.45 kB
├── assets/
│   ├── index-*.css        # 12.75 kB (3.33 kB gzipped)
│   └── index-*.js         # 245 kB (77.8 kB gzipped)
```

Go server serves from `web/dist/` in production.

## Next Steps

1. **Configure AWS Cognito** (if not already done)
   - Update Go server environment variables
   - Test sign up flow

2. **Add More Features**
   - Create item form
   - Edit/delete items
   - User profile editing

3. **Production Deployment**
   - Set up environment variables
   - Configure domain/SSL
   - Deploy Go server + built React app

## Support

- **React issues:** Check browser console + Network tab
- **API issues:** Check Go server logs
- **Build issues:** Run `npm run build` and check for errors
- **Type errors:** Run `npm run build` (TypeScript check included)
