# Frontend - React Vite SPA

This directory contains the React + TypeScript frontend for the AWS Go Server.

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool and dev server
- **CSS** - Styling (can be extended with Tailwind, etc.)

## Getting Started

### Install Dependencies

```bash
npm install
```

Or from the project root:

```bash
make frontend-install
```

### Development

Run the development server with hot reload:

```bash
npm run dev
```

Or from the project root:

```bash
make frontend-dev
```

The dev server will start at `http://localhost:5173` with:
- Hot Module Replacement (HMR)
- API proxy to Go backend at `http://localhost:8080`

**Important**: Make sure the Go backend is running before starting the frontend dev server.

### Build for Production

Build the optimized production bundle:

```bash
npm run build
```

Or from the project root:

```bash
make frontend-build
```

The build output will be in `web/dist/` directory.

## Development Workflow

### Option 1: Separate Dev Servers (Recommended for Development)

1. **Terminal 1** - Start the Go backend:
   ```bash
   make dev
   ```

2. **Terminal 2** - Start the React frontend:
   ```bash
   make frontend-dev
   ```

3. Open `http://localhost:5173` in your browser

**Benefits:**
- Fast hot reload for frontend changes
- API calls are proxied to the backend
- Best DX (Developer Experience)

### Option 2: Integrated Build (Production-like)

1. Build the frontend:
   ```bash
   make frontend-build
   ```

2. Start the Go server:
   ```bash
   make run
   ```

3. Open `http://localhost:8080` in your browser

**Benefits:**
- Tests the full production setup
- Single server serves both frontend and API
- Mimics production environment

## API Integration

### Proxy Configuration

In development, Vite proxies API requests to the Go backend:

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

This means:
- `http://localhost:5173/api/v1/items` → `http://localhost:8080/api/v1/items`
- `http://localhost:5173/healthz` → `http://localhost:8080/healthz`

### Making API Calls

Example API call in React:

```typescript
// Health check
const response = await fetch('/healthz')
const data = await response.json()

// API endpoint (requires authentication)
const response = await fetch('/api/v1/items', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
const items = await response.json()
```

## Project Structure

```
web/
├── public/           # Static assets
├── src/
│   ├── assets/      # Images, fonts, etc.
│   ├── App.tsx      # Main App component
│   ├── App.css      # App styles
│   ├── main.tsx     # Entry point
│   └── vite-env.d.ts # Vite type definitions
├── index.html       # HTML template
├── package.json     # Dependencies
├── tsconfig.json    # TypeScript config
└── vite.config.ts   # Vite config
```

## Environment Variables

Create a `.env` file in the `web/` directory for environment-specific config:

```bash
# Example .env
VITE_API_URL=http://localhost:8080
```

Access in code:

```typescript
const apiUrl = import.meta.env.VITE_API_URL
```

**Note**: All environment variables must be prefixed with `VITE_` to be exposed to the client.

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build locally
- `npm run lint` - Run ESLint (if configured)

## Adding Authentication

The backend uses AWS Cognito for authentication. To integrate with the frontend:

1. Install AWS Amplify or use fetch directly:
   ```bash
   npm install aws-amplify
   ```

2. Configure Cognito in your app:
   ```typescript
   import { Amplify } from 'aws-amplify'

   Amplify.configure({
     Auth: {
       region: 'us-east-1',
       userPoolId: 'your-user-pool-id',
       userPoolWebClientId: 'your-client-id'
     }
   })
   ```

3. Or use the backend auth endpoints:
   - `POST /api/v1/auth/signup`
   - `POST /api/v1/auth/login`
   - `POST /api/v1/auth/refresh`

## Next Steps

- Add routing with React Router
- Add state management (Zustand, Redux, etc.)
- Add UI component library (shadcn/ui, MUI, etc.)
- Add form validation (React Hook Form + Zod)
- Add API client library (TanStack Query)
- Set up testing (Vitest + React Testing Library)

## Deployment

The frontend is served by the Go backend in production:

1. Build the frontend: `make frontend-build`
2. Build the backend: `make build`
3. Deploy the Go binary (it will serve the frontend from `web/dist/`)

For separate deployment (CDN):
- Build and upload `web/dist/` to S3, CloudFront, Netlify, Vercel, etc.
- Update CORS settings in the Go backend to allow the CDN origin
