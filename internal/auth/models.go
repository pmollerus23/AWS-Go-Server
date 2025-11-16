package auth

import "time"

// User represents an authenticated user.
type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	IsAdmin  bool     `json:"is_admin"`
}

// Claims represents JWT token claims.
type Claims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	IsAdmin  bool     `json:"is_admin"`
	IssuedAt int64    `json:"iat"`
	ExpiresAt int64   `json:"exp"`
}

// TokenPair represents access and refresh tokens.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response.
type LoginResponse struct {
	User   User       `json:"user"`
	Tokens TokenPair  `json:"tokens"`
}

// Permission represents an authorization permission.
type Permission string

const (
	PermissionReadItems   Permission = "items:read"
	PermissionWriteItems  Permission = "items:write"
	PermissionDeleteItems Permission = "items:delete"
	PermissionAWSRead     Permission = "aws:read"
	PermissionAWSWrite    Permission = "aws:write"
	PermissionAdmin       Permission = "admin:*"
)

// Role represents a user role with permissions.
type Role struct {
	Name        string
	Permissions []Permission
}

// Predefined roles
var (
	RoleUser = Role{
		Name: "user",
		Permissions: []Permission{
			PermissionReadItems,
		},
	}

	RoleEditor = Role{
		Name: "editor",
		Permissions: []Permission{
			PermissionReadItems,
			PermissionWriteItems,
		},
	}

	RoleAdmin = Role{
		Name: "admin",
		Permissions: []Permission{
			PermissionReadItems,
			PermissionWriteItems,
			PermissionDeleteItems,
			PermissionAWSRead,
			PermissionAWSWrite,
			PermissionAdmin,
		},
	}
)

// GetRolePermissions returns permissions for a role name.
func GetRolePermissions(roleName string) []Permission {
	switch roleName {
	case "admin":
		return RoleAdmin.Permissions
	case "editor":
		return RoleEditor.Permissions
	case "user":
		return RoleUser.Permissions
	default:
		return []Permission{}
	}
}

// HasPermission checks if a user has a specific permission.
func (u *User) HasPermission(perm Permission) bool {
	// Admin has all permissions
	if u.IsAdmin {
		return true
	}

	// Check all roles
	for _, role := range u.Roles {
		permissions := GetRolePermissions(role)
		for _, p := range permissions {
			if p == perm || p == PermissionAdmin {
				return true
			}
		}
	}

	return false
}

// HasAnyRole checks if user has any of the specified roles.
func (u *User) HasAnyRole(roles ...string) bool {
	for _, userRole := range u.Roles {
		for _, requiredRole := range roles {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}
