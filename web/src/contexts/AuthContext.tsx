import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import type { PropsWithChildren, AuthState, User, LoginCredentials } from '../types';
import { authApi } from '../api';
import { storage, STORAGE_KEYS, getUserFromIdToken, isTokenExpired } from '../lib';

interface AuthContextValue extends AuthState {
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  updateUser: (user: Partial<User>) => void;
  refreshAuth: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const useAuth = (): AuthContextValue => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

interface AuthProviderProps extends PropsWithChildren {
  initialAuth?: Partial<AuthState>;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({
  children,
  initialAuth
}) => {
  const [authState, setAuthState] = useState<AuthState>({
    isAuthenticated: false,
    isLoading: true,
    user: null,
    error: null,
    ...initialAuth,
  });

  const login = useCallback(async (credentials: LoginCredentials): Promise<void> => {
    setAuthState(prev => ({ ...prev, isLoading: true, error: null }));

    try {
      const response = await authApi.login(credentials);

      // Store tokens
      storage.set(STORAGE_KEYS.ACCESS_TOKEN, response.tokens.access_token);
      storage.set(STORAGE_KEYS.ID_TOKEN, response.tokens.id_token);
      if (response.tokens.refresh_token) {
        storage.set(STORAGE_KEYS.REFRESH_TOKEN, response.tokens.refresh_token);
      }
      storage.set(STORAGE_KEYS.USER_EMAIL, credentials.email);

      // Extract user from ID token
      const user = getUserFromIdToken(response.tokens.id_token);

      if (!user) {
        throw new Error('Failed to extract user from token');
      }

      setAuthState({
        isAuthenticated: true,
        isLoading: false,
        user,
        error: null,
      });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Login failed';
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage,
      }));
      throw error;
    }
  }, []);

  const logout = useCallback(async (): Promise<void> => {
    try {
      // Clear tokens from storage
      storage.remove(STORAGE_KEYS.ACCESS_TOKEN);
      storage.remove(STORAGE_KEYS.ID_TOKEN);
      storage.remove(STORAGE_KEYS.REFRESH_TOKEN);
      storage.remove(STORAGE_KEYS.USER_EMAIL);

      setAuthState({
        isAuthenticated: false,
        isLoading: false,
        user: null,
        error: null,
      });
    } catch (error) {
      console.error('Logout error:', error);
      throw error;
    }
  }, []);

  const updateUser = useCallback((userData: Partial<User>): void => {
    setAuthState(prev => ({
      ...prev,
      user: prev.user ? { ...prev.user, ...userData } : null,
    }));
  }, []);

  const refreshAuth = useCallback(async (): Promise<void> => {
    setAuthState(prev => ({ ...prev, isLoading: true }));

    try {
      const accessToken = storage.get<string>(STORAGE_KEYS.ACCESS_TOKEN);
      const idToken = storage.get<string>(STORAGE_KEYS.ID_TOKEN);
      const refreshToken = storage.get<string>(STORAGE_KEYS.REFRESH_TOKEN);
      const userEmail = storage.get<string>(STORAGE_KEYS.USER_EMAIL);

      // No tokens found - user is not authenticated
      if (!accessToken || !idToken) {
        setAuthState({
          isAuthenticated: false,
          isLoading: false,
          user: null,
          error: null,
        });
        return;
      }

      // Check if ID token is expired
      if (isTokenExpired(idToken)) {
        // Try to refresh if we have a refresh token
        if (refreshToken && userEmail) {
          try {
            const response = await authApi.refreshToken(refreshToken, userEmail);

            // Update stored tokens
            storage.set(STORAGE_KEYS.ACCESS_TOKEN, response.tokens.access_token);
            storage.set(STORAGE_KEYS.ID_TOKEN, response.tokens.id_token);

            const user = getUserFromIdToken(response.tokens.id_token);

            if (!user) {
              throw new Error('Failed to extract user from refreshed token');
            }

            setAuthState({
              isAuthenticated: true,
              isLoading: false,
              user,
              error: null,
            });
            return;
          } catch (error) {
            // Refresh failed - clear auth
            console.error('Token refresh failed:', error);
            // Clear tokens directly instead of calling logout
            storage.remove(STORAGE_KEYS.ACCESS_TOKEN);
            storage.remove(STORAGE_KEYS.ID_TOKEN);
            storage.remove(STORAGE_KEYS.REFRESH_TOKEN);
            storage.remove(STORAGE_KEYS.USER_EMAIL);
            setAuthState({
              isAuthenticated: false,
              isLoading: false,
              user: null,
              error: null,
            });
            return;
          }
        } else {
          // No refresh token - clear auth
          storage.remove(STORAGE_KEYS.ACCESS_TOKEN);
          storage.remove(STORAGE_KEYS.ID_TOKEN);
          storage.remove(STORAGE_KEYS.REFRESH_TOKEN);
          storage.remove(STORAGE_KEYS.USER_EMAIL);
          setAuthState({
            isAuthenticated: false,
            isLoading: false,
            user: null,
            error: null,
          });
          return;
        }
      }

      // Token is still valid - extract user
      const user = getUserFromIdToken(idToken);

      if (!user) {
        throw new Error('Failed to extract user from token');
      }

      setAuthState({
        isAuthenticated: true,
        isLoading: false,
        user,
        error: null,
      });
    } catch (error) {
      console.error('Auth refresh error:', error);
      setAuthState({
        isAuthenticated: false,
        isLoading: false,
        user: null,
        error: error instanceof Error ? error.message : 'Auth refresh failed',
      });
    }
  }, []); // No dependencies - stable function

  useEffect(() => {
    // Initialize auth state on mount only
    refreshAuth();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Only run once on mount

  const value: AuthContextValue = {
    ...authState,
    login,
    logout,
    updateUser,
    refreshAuth,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
