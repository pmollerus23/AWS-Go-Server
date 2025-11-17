# React Microfrontend Architecture

This document outlines the architecture and folder structure of the React microfrontend application.

## Table of Contents

- [Folder Structure](#folder-structure)
- [Architecture Patterns](#architecture-patterns)
- [State Management](#state-management)
- [API Integration](#api-integration)
- [Component Guidelines](#component-guidelines)

## Folder Structure

```
src/
├── api/              # API client and service modules
│   ├── client.ts         # Base API client with fetch wrapper
│   ├── auth.api.ts       # Authentication API endpoints
│   └── index.ts          # API exports
│
├── components/       # Feature/business components
│   ├── UserCard.tsx      # Example feature component
│   ├── ProtectedRoute.tsx # Route protection component
│   ├── ErrorBoundary.tsx  # Error handling component
│   └── index.ts          # Component exports
│
├── contexts/         # React Context providers
│   ├── AuthContext.tsx   # Authentication context & provider
│   └── index.ts          # Context exports
│
├── hooks/            # Custom React hooks
│   ├── useQuery.ts       # Data fetching hook
│   ├── useMutation.ts    # Data mutation hook
│   └── index.ts          # Hook exports
│
├── lib/              # Utility functions and helpers
│   ├── constants.ts      # Application constants
│   ├── utils.ts          # General utility functions
│   ├── storage.ts        # LocalStorage wrapper
│   ├── validators.ts     # Input validation functions
│   └── index.ts          # Library exports
│
├── pages/            # Page components (route-level)
│   ├── HomePage.tsx      # Home/Dashboard page
│   ├── ProfilePage.tsx   # User profile page
│   ├── LoginPage.tsx     # Login page
│   └── index.ts          # Page exports
│
├── shell/            # Application shell/layout
│   ├── Shell.tsx         # Main shell wrapper
│   ├── Header.tsx        # Header component
│   ├── Sidebar.tsx       # Sidebar navigation
│   ├── Footer.tsx        # Footer component
│   └── index.ts          # Shell exports
│
├── types/            # TypeScript type definitions
│   ├── auth.types.ts     # Authentication types
│   ├── api.types.ts      # API-related types
│   ├── common.types.ts   # Common utility types
│   └── index.ts          # Type exports
│
├── ui/               # Reusable UI components (design system)
│   ├── Button.tsx        # Button component
│   ├── Card.tsx          # Card component
│   ├── Input.tsx         # Input component
│   ├── Modal.tsx         # Modal component
│   ├── Spinner.tsx       # Loading spinner
│   └── index.ts          # UI exports
│
├── styles/           # Global styles
│   └── globals.css       # Global CSS variables and base styles
│
├── App.tsx           # Root application component
├── main.tsx          # Application entry point
└── vite-env.d.ts     # Vite environment types
```

## Architecture Patterns

### 1. Microfrontend Shell Pattern

The application uses a shell architecture where:

- **Shell** (`src/shell/`): Provides the main layout structure (Header, Sidebar, Footer)
- **Pages** (`src/pages/`): Route-level components that render within the shell
- **Components** (`src/components/`): Feature-specific components used within pages

### 2. Component Composition

All components follow these principles:

- **Functional Components**: Only functional components with hooks
- **TypeScript**: Strongly typed with interfaces and utility types
- **Props Interface**: Each component has a dedicated props interface
- **Composition over Inheritance**: Components are composed together

Example:

```typescript
interface ButtonProps extends BaseComponentProps, PropsWithChildren {
  variant?: 'primary' | 'secondary' | 'danger';
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
}

export const Button: React.FC<ButtonProps> = ({ children, variant = 'primary', onClick }) => {
  // Component implementation
};
```

### 3. Separation of Concerns

- **UI Components** (`src/ui/`): Pure presentational components, no business logic
- **Feature Components** (`src/components/`): Business logic and feature-specific behavior
- **Pages** (`src/pages/`): Route-level orchestration of components
- **Hooks** (`src/hooks/`): Reusable stateful logic
- **API Layer** (`src/api/`): All server communication

## State Management

### Context for Global State

Authentication and user profile are managed via React Context:

```typescript
// Usage
const { user, isAuthenticated, login, logout } = useAuth();
```

The `AuthContext` provides:
- `isAuthenticated`: Boolean authentication status
- `user`: Current user profile
- `login()`: Login method
- `logout()`: Logout method
- `updateUser()`: Update user data
- `refreshAuth()`: Refresh authentication state

### Parent-Child State for Local State

Shared state between parent and child components uses standard React state:

```typescript
const [sharedData, setSharedData] = useState<Data>();

<Parent>
  <Child data={sharedData} onUpdate={setSharedData} />
</Parent>
```

## API Integration

### useQuery Hook

All API GET requests use the `useQuery` hook:

```typescript
const { data, isLoading, isError, error, refetch } = useQuery(
  async () => apiClient.get('/endpoint'),
  {
    enabled: true,
    refetchOnMount: true,
    refetchInterval: 30000, // 30 seconds
    onSuccess: (data) => console.log('Success:', data),
    onError: (error) => console.error('Error:', error),
  }
);
```

Features:
- Automatic loading and error states
- Retry logic with configurable attempts
- Refetch on mount
- Polling with refetchInterval
- Success/error callbacks

### useMutation Hook

All API POST/PUT/PATCH/DELETE requests use the `useMutation` hook:

```typescript
const mutation = useMutation(
  async (variables) => apiClient.post('/endpoint', variables),
  {
    onSuccess: (data) => console.log('Success:', data),
    onError: (error) => console.error('Error:', error),
  }
);

// Trigger mutation
mutation.mutate({ name: 'value' });
// or await
const result = await mutation.mutateAsync({ name: 'value' });
```

### API Client

The API client (`src/api/client.ts`) provides:

- Automatic authentication header injection
- Request/response interceptors
- Error handling with typed errors
- Query parameter serialization

```typescript
// GET request
const data = await apiClient.get<User>('/users/123');

// POST request
const created = await apiClient.post<User>('/users', { name: 'John' });

// PUT request
const updated = await apiClient.put<User>('/users/123', { name: 'Jane' });

// DELETE request
await apiClient.delete('/users/123');
```

## Component Guidelines

### TypeScript Types

Use utility types from `src/types/common.types.ts`:

```typescript
// Props with children
interface MyProps extends PropsWithChildren {
  title: string;
}

// Base component props (className, id, testId)
interface MyProps extends BaseComponentProps {
  variant: 'primary' | 'secondary';
}

// Nullable types
type MaybeUser = Nullable<User>;

// Async functions
const fetchData: AsyncFunction<User[]> = async () => {
  return await apiClient.get('/users');
};
```

### Component Structure

```typescript
import type { PropsWithChildren, BaseComponentProps } from '../types';

interface MyComponentProps extends BaseComponentProps, PropsWithChildren {
  title: string;
  onAction?: () => void;
}

export const MyComponent: React.FC<MyComponentProps> = ({
  title,
  onAction,
  children,
  className,
  id,
  testId,
}) => {
  return (
    <div className={className} id={id} data-testid={testId}>
      <h2>{title}</h2>
      {children}
      {onAction && <button onClick={onAction}>Action</button>}
    </div>
  );
};
```

### Error Handling

Wrap your app with `ErrorBoundary`:

```typescript
<ErrorBoundary
  fallback={<ErrorPage />}
  onError={(error, errorInfo) => {
    logErrorToService(error, errorInfo);
  }}
>
  <App />
</ErrorBoundary>
```

### Protected Routes

Wrap authenticated pages with `ProtectedRoute`:

```typescript
<ProtectedRoute redirectTo="/login" requiredRoles={['admin']}>
  <AdminPage />
</ProtectedRoute>
```

## Environment Variables

Create a `.env` file based on `.env.example`:

```env
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=My App
```

Access in code:

```typescript
const apiUrl = import.meta.env.VITE_API_URL;
```

## Best Practices

1. **Always use TypeScript**: Define interfaces for all props and data structures
2. **Component naming**: Use PascalCase for components, camelCase for functions/variables
3. **File organization**: One component per file, export from index.ts
4. **Avoid prop drilling**: Use Context for deeply nested shared state
5. **Error handling**: Always handle loading and error states in data fetching
6. **Accessibility**: Include ARIA labels and semantic HTML
7. **Testing**: Use `testId` props for reliable test selectors

## Next Steps

1. Add routing library (React Router, TanStack Router, etc.)
2. Implement additional API endpoints in `src/api/`
3. Create more reusable UI components in `src/ui/`
4. Add form validation library integration
5. Set up unit and integration tests
6. Configure CI/CD pipeline
