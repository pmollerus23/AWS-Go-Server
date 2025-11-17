import { Navigate } from 'react-router-dom';
import { useAuth } from '../contexts';
import type { PropsWithChildren } from '../types';

interface ProtectedRouteProps extends PropsWithChildren {
  redirectTo?: string;
  requiredRoles?: string[];
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  redirectTo = '/login',
  requiredRoles = [],
}) => {
  const { isAuthenticated, isLoading, user } = useAuth();

  if (isLoading) {
    return (
      <div className="protected-route-loading">
        <p>Loading...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to={redirectTo} replace />;
  }

  // Check role requirements
  if (requiredRoles.length > 0) {
    const userGroups = user?.groups || [];
    const hasRequiredRole = requiredRoles.some(role => userGroups.includes(role));

    if (!hasRequiredRole) {
      return (
        <div className="protected-route-forbidden">
          <h1>Access Denied</h1>
          <p>You don't have permission to access this page.</p>
        </div>
      );
    }
  }

  return <>{children}</>;
};
