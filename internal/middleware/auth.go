package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/pmollerus23/go-aws-server/internal/auth"
)

// AuthService defines the interface for authentication services.
type AuthService interface {
	ValidateToken(ctx context.Context, token string) (*auth.Claims, error)
}

// Authenticate is middleware that validates JWT tokens from AWS Cognito.
func Authenticate(authService AuthService, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("missing authorization header",
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
				return
			}

			// Check for Bearer token format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				logger.Warn("invalid authorization header format",
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, "Unauthorized: invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(r.Context(), token)
			if err != nil {
				logger.Warn("token validation failed",
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			// Convert claims to user
			user := &auth.User{
				ID:       claims.UserID,
				Email:    claims.Email,
				Username: claims.Username,
				Roles:    claims.Roles,
				IsAdmin:  claims.IsAdmin,
			}

			// Add user to context
			ctx := auth.WithUser(r.Context(), user)

			logger.Info("request authenticated",
				"user_id", user.ID,
				"email", user.Email,
				"path", r.URL.Path,
				"method", r.Method,
			)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePermission is middleware that checks if the authenticated user has a specific permission.
func RequirePermission(permission auth.Permission, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.GetUser(r.Context())
			if err != nil {
				logger.Warn("no user in context for permission check",
					"permission", permission,
					"path", r.URL.Path,
				)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !user.HasPermission(permission) {
				logger.Warn("user lacks required permission",
					"user_id", user.ID,
					"permission", permission,
					"path", r.URL.Path,
				)
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole is middleware that checks if the authenticated user has any of the specified roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.GetUser(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !user.HasAnyRole(roles...) {
				http.Error(w, "Forbidden: insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin is middleware that checks if the authenticated user is an admin.
func RequireAdmin(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := auth.GetUser(r.Context())
			if err != nil {
				logger.Warn("no user in context for admin check",
					"path", r.URL.Path,
				)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !user.IsAdmin {
				logger.Warn("non-admin user attempted admin access",
					"user_id", user.ID,
					"path", r.URL.Path,
				)
				http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
