// JWT Claims from Cognito ID token
export interface CognitoUser {
  sub: string; // User ID
  email: string;
  email_verified?: boolean;
  name?: string;
  given_name?: string;
  family_name?: string;
  'cognito:username'?: string;
  'cognito:groups'?: string[];
  exp?: number;
  iat?: number;
}

export interface User {
  id: string;
  email: string;
  username: string;
  name?: string;
  firstName?: string;
  lastName?: string;
  emailVerified?: boolean;
  groups?: string[];
}

export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: User | null;
  error: string | null;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignUpData {
  email: string;
  password: string;
  name?: string;
}

export interface ConfirmSignUpData {
  email: string;
  code: string;
}

export interface ForgotPasswordData {
  email: string;
}

export interface ResetPasswordData {
  email: string;
  code: string;
  newPassword: string;
}

// Cognito Tokens from Go server
export interface CognitoTokens {
  access_token: string;
  id_token: string;
  refresh_token?: string;
  expires_in: number;
  token_type: string;
}

// API Response wrappers
export interface SignUpResponse {
  message: string;
  email: string;
}

export interface ConfirmSignUpResponse {
  message: string;
}

export interface LoginResponse {
  message: string;
  tokens: CognitoTokens;
}

export interface RefreshTokenResponse {
  message: string;
  tokens: CognitoTokens;
}

export interface ForgotPasswordResponse {
  message: string;
}

export interface ResetPasswordResponse {
  message: string;
}
