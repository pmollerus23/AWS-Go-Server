package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pmollerus23/go-aws-server/internal/config"
)

var (
	ErrInvalidCredentials   = errors.New("invalid email or password")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotConfirmed     = errors.New("user email not verified")
	ErrInvalidVerification  = errors.New("invalid verification code")
	ErrPasswordResetRequired = errors.New("password reset required")
)

// CognitoService handles AWS Cognito authentication operations.
type CognitoService struct {
	client       *cognito.Client
	cfg          config.CognitoConfig
	logger       *slog.Logger
	jwksCache    jwk.Set
	jwksURL      string
	cacheExpiry  time.Time
}

// NewCognitoService creates a new Cognito service.
func NewCognitoService(client *cognito.Client, cfg config.CognitoConfig, logger *slog.Logger) *CognitoService {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		cfg.Region, cfg.UserPoolID)

	return &CognitoService{
		client:  client,
		cfg:     cfg,
		logger:  logger,
		jwksURL: jwksURL,
	}
}

// SignUp registers a new user with Cognito.
func (s *CognitoService) SignUp(ctx context.Context, email, password, name string) error {
	secretHash := s.calculateSecretHash(email)

	input := &cognito.SignUpInput{
		ClientId:   aws.String(s.cfg.ClientID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(email),
		Password:   aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	if name != "" {
		input.UserAttributes = append(input.UserAttributes, types.AttributeType{
			Name:  aws.String("name"),
			Value: aws.String(name),
		})
	}

	_, err := s.client.SignUp(ctx, input)
	if err != nil {
		var usernameExists *types.UsernameExistsException
		if errors.As(err, &usernameExists) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("cognito signup failed: %w", err)
	}

	s.logger.Info("user signed up successfully", "email", email)
	return nil
}

// ConfirmSignUp confirms a user's signup with the verification code sent via email.
func (s *CognitoService) ConfirmSignUp(ctx context.Context, email, code string) error {
	secretHash := s.calculateSecretHash(email)

	input := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(s.cfg.ClientID),
		SecretHash:       aws.String(secretHash),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
	}

	_, err := s.client.ConfirmSignUp(ctx, input)
	if err != nil {
		var codeExpired *types.ExpiredCodeException
		var codeMismatch *types.CodeMismatchException
		if errors.As(err, &codeExpired) || errors.As(err, &codeMismatch) {
			return ErrInvalidVerification
		}
		return fmt.Errorf("cognito confirm signup failed: %w", err)
	}

	s.logger.Info("user confirmed successfully", "email", email)
	return nil
}

// Login authenticates a user and returns JWT tokens.
func (s *CognitoService) Login(ctx context.Context, email, password string) (*CognitoTokens, error) {
	secretHash := s.calculateSecretHash(email)

	input := &cognito.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(s.cfg.ClientID),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	}

	result, err := s.client.InitiateAuth(ctx, input)
	if err != nil {
		var notAuthorized *types.NotAuthorizedException
		var userNotConfirmed *types.UserNotConfirmedException
		var passwordReset *types.PasswordResetRequiredException

		if errors.As(err, &notAuthorized) {
			return nil, ErrInvalidCredentials
		}
		if errors.As(err, &userNotConfirmed) {
			return nil, ErrUserNotConfirmed
		}
		if errors.As(err, &passwordReset) {
			return nil, ErrPasswordResetRequired
		}

		return nil, fmt.Errorf("cognito login failed: %w", err)
	}

	if result.AuthenticationResult == nil {
		return nil, fmt.Errorf("authentication result is nil")
	}

	tokens := &CognitoTokens{
		AccessToken:  aws.ToString(result.AuthenticationResult.AccessToken),
		IDToken:      aws.ToString(result.AuthenticationResult.IdToken),
		RefreshToken: aws.ToString(result.AuthenticationResult.RefreshToken),
		ExpiresIn:    result.AuthenticationResult.ExpiresIn,
		TokenType:    aws.ToString(result.AuthenticationResult.TokenType),
	}

	s.logger.Info("user logged in successfully", "email", email)
	return tokens, nil
}

// RefreshToken refreshes access and ID tokens using a refresh token.
func (s *CognitoService) RefreshToken(ctx context.Context, refreshToken, email string) (*CognitoTokens, error) {
	secretHash := s.calculateSecretHash(email)

	input := &cognito.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(s.cfg.ClientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   secretHash,
		},
	}

	result, err := s.client.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("cognito refresh token failed: %w", err)
	}

	if result.AuthenticationResult == nil {
		return nil, fmt.Errorf("authentication result is nil")
	}

	tokens := &CognitoTokens{
		AccessToken: aws.ToString(result.AuthenticationResult.AccessToken),
		IDToken:     aws.ToString(result.AuthenticationResult.IdToken),
		ExpiresIn:   result.AuthenticationResult.ExpiresIn,
		TokenType:   aws.ToString(result.AuthenticationResult.TokenType),
	}

	s.logger.Info("token refreshed successfully")
	return tokens, nil
}

// ValidateToken validates a JWT token from Cognito using JWKS.
func (s *CognitoService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	// Refresh JWKS cache if expired
	if err := s.refreshJWKSCache(ctx); err != nil {
		return nil, fmt.Errorf("failed to refresh JWKS cache: %w", err)
	}

	// Parse and validate token
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(s.jwksCache),
		jwt.WithValidate(true),
	)
	if err != nil {
		s.logger.Error("token validation failed", "error", err)
		return nil, ErrInvalidToken
	}

	// Verify issuer
	expectedIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s",
		s.cfg.Region, s.cfg.UserPoolID)
	if token.Issuer() != expectedIssuer {
		return nil, ErrInvalidToken
	}

	// Verify token use (should be "access" for access tokens)
	tokenUse, ok := token.Get("token_use")
	if !ok || tokenUse != "access" {
		return nil, ErrInvalidToken
	}

	// Extract claims
	claims := &Claims{
		UserID:    token.Subject(),
		ExpiresAt: token.Expiration().Unix(),
		IssuedAt:  token.IssuedAt().Unix(),
	}

	// Extract cognito:username
	if username, ok := token.Get("cognito:username"); ok {
		if usernameStr, ok := username.(string); ok {
			claims.Username = usernameStr
		}
	}

	// Extract email
	if email, ok := token.Get("email"); ok {
		if emailStr, ok := email.(string); ok {
			claims.Email = emailStr
		}
	}

	// Extract cognito:groups (roles)
	if groups, ok := token.Get("cognito:groups"); ok {
		if groupsSlice, ok := groups.([]interface{}); ok {
			roles := make([]string, 0, len(groupsSlice))
			for _, g := range groupsSlice {
				if role, ok := g.(string); ok {
					roles = append(roles, role)
				}
			}
			claims.Roles = roles
		}
	}

	// Check if user is admin (has "admin" group)
	for _, role := range claims.Roles {
		if role == "admin" {
			claims.IsAdmin = true
			break
		}
	}

	return claims, nil
}

// ForgotPassword initiates the forgot password flow.
func (s *CognitoService) ForgotPassword(ctx context.Context, email string) error {
	secretHash := s.calculateSecretHash(email)

	input := &cognito.ForgotPasswordInput{
		ClientId:   aws.String(s.cfg.ClientID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(email),
	}

	_, err := s.client.ForgotPassword(ctx, input)
	if err != nil {
		return fmt.Errorf("cognito forgot password failed: %w", err)
	}

	s.logger.Info("forgot password initiated", "email", email)
	return nil
}

// ConfirmForgotPassword confirms password reset with the code.
func (s *CognitoService) ConfirmForgotPassword(ctx context.Context, email, code, newPassword string) error {
	secretHash := s.calculateSecretHash(email)

	input := &cognito.ConfirmForgotPasswordInput{
		ClientId:         aws.String(s.cfg.ClientID),
		SecretHash:       aws.String(secretHash),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(code),
		Password:         aws.String(newPassword),
	}

	_, err := s.client.ConfirmForgotPassword(ctx, input)
	if err != nil {
		var codeExpired *types.ExpiredCodeException
		var codeMismatch *types.CodeMismatchException
		if errors.As(err, &codeExpired) || errors.As(err, &codeMismatch) {
			return ErrInvalidVerification
		}
		return fmt.Errorf("cognito confirm forgot password failed: %w", err)
	}

	s.logger.Info("password reset successfully", "email", email)
	return nil
}

// calculateSecretHash calculates the secret hash required for Cognito API calls.
func (s *CognitoService) calculateSecretHash(username string) string {
	message := username + s.cfg.ClientID
	key := []byte(s.cfg.ClientSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// refreshJWKSCache refreshes the JWKS cache if it's expired or not yet loaded.
func (s *CognitoService) refreshJWKSCache(ctx context.Context) error {
	// Check if cache is still valid
	if s.jwksCache != nil && time.Now().Before(s.cacheExpiry) {
		return nil
	}

	// Fetch JWKS
	cache := jwk.NewCache(ctx)
	if err := cache.Register(s.jwksURL); err != nil {
		return fmt.Errorf("failed to register JWKS URL: %w", err)
	}

	cached, err := cache.Refresh(ctx, s.jwksURL)
	if err != nil {
		return fmt.Errorf("failed to refresh JWKS: %w", err)
	}

	s.jwksCache = cached
	s.cacheExpiry = time.Now().Add(1 * time.Hour) // Cache for 1 hour

	s.logger.Info("JWKS cache refreshed")
	return nil
}

// CognitoTokens represents tokens returned from Cognito authentication.
type CognitoTokens struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int32  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
