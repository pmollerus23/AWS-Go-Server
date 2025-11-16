package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
)

// JWTService handles JWT token creation and validation.
type JWTService struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// NewJWTService creates a new JWT service.
func NewJWTService(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:            []byte(secretKey),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		issuer:               "aws-go-server",
	}
}

// GenerateTokenPair generates access and refresh tokens for a user.
func (s *JWTService) GenerateTokenPair(user *User) (*TokenPair, error) {
	// Generate access token
	accessToken, expiresAt, err := s.generateToken(user, s.accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, _, err := s.generateRefreshToken(user.ID, s.refreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}

// generateToken creates a JWT access token.
func (s *JWTService) generateToken(user *User, duration time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"roles":    user.Roles,
		"is_admin": user.IsAdmin,
		"iat":      now.Unix(),
		"exp":      expiresAt.Unix(),
		"iss":      s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateRefreshToken creates a refresh token.
func (s *JWTService) generateRefreshToken(userID string, duration time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
		"iss":     s.issuer,
		"jti":     generateJTI(), // Unique token ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Extract and validate claims
	userClaims := &Claims{}

	if userID, ok := claims["user_id"].(string); ok {
		userClaims.UserID = userID
	} else {
		return nil, ErrInvalidToken
	}

	if email, ok := claims["email"].(string); ok {
		userClaims.Email = email
	}

	if username, ok := claims["username"].(string); ok {
		userClaims.Username = username
	}

	if isAdmin, ok := claims["is_admin"].(bool); ok {
		userClaims.IsAdmin = isAdmin
	}

	// Extract roles
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		roles := make([]string, 0, len(rolesInterface))
		for _, r := range rolesInterface {
			if roleStr, ok := r.(string); ok {
				roles = append(roles, roleStr)
			}
		}
		userClaims.Roles = roles
	}

	if iat, ok := claims["iat"].(float64); ok {
		userClaims.IssuedAt = int64(iat)
	}

	if exp, ok := claims["exp"].(float64); ok {
		userClaims.ExpiresAt = int64(exp)
	}

	return userClaims, nil
}

// ClaimsToUser converts JWT claims to a User object.
func (s *JWTService) ClaimsToUser(claims *Claims) *User {
	return &User{
		ID:       claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
		Roles:    claims.Roles,
		IsAdmin:  claims.IsAdmin,
	}
}

// generateJTI generates a unique JWT ID.
func generateJTI() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
