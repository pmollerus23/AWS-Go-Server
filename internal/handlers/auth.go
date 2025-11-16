package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/pmollerus23/go-aws-server/internal/auth"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	SignUp(ctx context.Context, email, password, name string) error
	ConfirmSignUp(ctx context.Context, email, code string) error
	Login(ctx context.Context, email, password string) (*auth.CognitoTokens, error)
	RefreshToken(ctx context.Context, refreshToken, email string) (*auth.CognitoTokens, error)
	ForgotPassword(ctx context.Context, email string) error
	ConfirmForgotPassword(ctx context.Context, email, code, newPassword string) error
}

// SignUpRequest represents the signup request payload.
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

// Valid validates the signup request.
func (r SignUpRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.Email == "" {
		problems["email"] = "email is required"
	}
	if r.Password == "" {
		problems["password"] = "password is required"
	}
	if len(r.Password) < 8 {
		problems["password"] = "password must be at least 8 characters"
	}

	return problems
}

// SignUpResponse represents the signup response.
type SignUpResponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

// HandleSignUp handles user registration.
//
//	@Summary		Sign up a new user
//	@Description	Register a new user account with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SignUpRequest	true	"Signup request"
//	@Success		201		{object}	SignUpResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		409		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/v1/auth/signup [post]
func HandleSignUp(logger *slog.Logger, authService AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[SignUpRequest](r)
		if err != nil {
			logger.Error("failed to decode signup request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = authService.SignUp(r.Context(), req.Email, req.Password, req.Name)
		if err != nil {
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				encode(w, r, http.StatusConflict, map[string]interface{}{
					"error": "user already exists",
				})
				return
			}
			logger.Error("signup failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := SignUpResponse{
			Message: "User registered successfully. Please check your email for verification code.",
			Email:   req.Email,
		}

		encode(w, r, http.StatusCreated, resp)
	})
}

// ConfirmSignUpRequest represents the confirm signup request.
type ConfirmSignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// Valid validates the confirm signup request.
func (r ConfirmSignUpRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.Email == "" {
		problems["email"] = "email is required"
	}
	if r.Code == "" {
		problems["code"] = "verification code is required"
	}

	return problems
}

// ConfirmSignUpResponse represents the confirm signup response.
type ConfirmSignUpResponse struct {
	Message string `json:"message"`
}

// HandleConfirmSignUp handles email verification.
//
//	@Summary		Confirm user signup
//	@Description	Verify user email with confirmation code
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ConfirmSignUpRequest	true	"Confirmation request"
//	@Success		200		{object}	ConfirmSignUpResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/v1/auth/confirm [post]
func HandleConfirmSignUp(logger *slog.Logger, authService AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[ConfirmSignUpRequest](r)
		if err != nil {
			logger.Error("failed to decode confirm request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = authService.ConfirmSignUp(r.Context(), req.Email, req.Code)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidVerification) {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error": "invalid or expired verification code",
				})
				return
			}
			logger.Error("confirm signup failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := ConfirmSignUpResponse{
			Message: "Email verified successfully. You can now login.",
		}

		encode(w, r, http.StatusOK, resp)
	})
}

// LoginRequest represents the login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Valid validates the login request.
func (r LoginRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.Email == "" {
		problems["email"] = "email is required"
	}
	if r.Password == "" {
		problems["password"] = "password is required"
	}

	return problems
}

// LoginResponse represents the login response.
type LoginResponse struct {
	Message string               `json:"message"`
	Tokens  *auth.CognitoTokens `json:"tokens"`
}

// HandleLogin handles user authentication.
//
//	@Summary		Login
//	@Description	Authenticate user and receive JWT tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/v1/auth/login [post]
func HandleLogin(logger *slog.Logger, authService AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[LoginRequest](r)
		if err != nil {
			logger.Error("failed to decode login request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tokens, err := authService.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				encode(w, r, http.StatusUnauthorized, map[string]interface{}{
					"error": "invalid email or password",
				})
				return
			}
			if errors.Is(err, auth.ErrUserNotConfirmed) {
				encode(w, r, http.StatusUnauthorized, map[string]interface{}{
					"error": "email not verified. Please check your email for verification code.",
				})
				return
			}
			logger.Error("login failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := LoginResponse{
			Message: "Login successful",
			Tokens:  tokens,
		}

		encode(w, r, http.StatusOK, resp)
	})
}

// RefreshTokenRequest represents the refresh token request.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
}

// Valid validates the refresh token request.
func (r RefreshTokenRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.RefreshToken == "" {
		problems["refresh_token"] = "refresh token is required"
	}
	if r.Email == "" {
		problems["email"] = "email is required"
	}

	return problems
}

// RefreshTokenResponse represents the refresh token response.
type RefreshTokenResponse struct {
	Message string               `json:"message"`
	Tokens  *auth.CognitoTokens `json:"tokens"`
}

// HandleRefreshToken handles token refresh.
//
//	@Summary		Refresh tokens
//	@Description	Refresh access and ID tokens using refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RefreshTokenRequest	true	"Refresh token request"
//	@Success		200		{object}	RefreshTokenResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/v1/auth/refresh [post]
func HandleRefreshToken(logger *slog.Logger, authService AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[RefreshTokenRequest](r)
		if err != nil {
			logger.Error("failed to decode refresh request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tokens, err := authService.RefreshToken(r.Context(), req.RefreshToken, req.Email)
		if err != nil {
			logger.Error("token refresh failed", "error", err)
			encode(w, r, http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid refresh token",
			})
			return
		}

		resp := RefreshTokenResponse{
			Message: "Tokens refreshed successfully",
			Tokens:  tokens,
		}

		encode(w, r, http.StatusOK, resp)
	})
}

// ForgotPasswordRequest represents the forgot password request.
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// Valid validates the forgot password request.
func (r ForgotPasswordRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.Email == "" {
		problems["email"] = "email is required"
	}

	return problems
}

// ForgotPasswordResponse represents the forgot password response.
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// HandleForgotPassword handles password reset request.
//
//	@Summary		Forgot password
//	@Description	Request password reset code via email
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ForgotPasswordRequest	true	"Forgot password request"
//	@Success		200		{object}	ForgotPasswordResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/v1/auth/forgot-password [post]
func HandleForgotPassword(logger *slog.Logger, authService AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[ForgotPasswordRequest](r)
		if err != nil {
			logger.Error("failed to decode forgot password request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = authService.ForgotPassword(r.Context(), req.Email)
		if err != nil {
			logger.Error("forgot password failed", "error", err)
			// Don't reveal if user exists or not
		}

		resp := ForgotPasswordResponse{
			Message: "If the email exists, a password reset code has been sent.",
		}

		encode(w, r, http.StatusOK, resp)
	})
}

// ConfirmForgotPasswordRequest represents the confirm forgot password request.
type ConfirmForgotPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

// Valid validates the confirm forgot password request.
func (r ConfirmForgotPasswordRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if r.Email == "" {
		problems["email"] = "email is required"
	}
	if r.Code == "" {
		problems["code"] = "verification code is required"
	}
	if r.NewPassword == "" {
		problems["new_password"] = "new password is required"
	}
	if len(r.NewPassword) < 8 {
		problems["new_password"] = "password must be at least 8 characters"
	}

	return problems
}

// ConfirmForgotPasswordResponse represents the confirm forgot password response.
type ConfirmForgotPasswordResponse struct {
	Message string `json:"message"`
}

// HandleConfirmForgotPassword handles password reset confirmation.
//
//	@Summary		Reset password
//	@Description	Reset password with verification code
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ConfirmForgotPasswordRequest	true	"Reset password request"
//	@Success		200		{object}	ConfirmForgotPasswordResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/v1/auth/reset-password [post]
func HandleConfirmForgotPassword(logger *slog.Logger, authService AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[ConfirmForgotPasswordRequest](r)
		if err != nil {
			logger.Error("failed to decode reset password request", "error", err)
			if len(problems) > 0 {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error":    "validation failed",
					"problems": problems,
				})
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		err = authService.ConfirmForgotPassword(r.Context(), req.Email, req.Code, req.NewPassword)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidVerification) {
				encode(w, r, http.StatusBadRequest, map[string]interface{}{
					"error": "invalid or expired verification code",
				})
				return
			}
			logger.Error("reset password failed", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := ConfirmForgotPasswordResponse{
			Message: "Password reset successfully. You can now login with your new password.",
		}

		encode(w, r, http.StatusOK, resp)
	})
}
