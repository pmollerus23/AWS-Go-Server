import { apiClient } from './client';
import type {
  LoginCredentials,
  SignUpData,
  ConfirmSignUpData,
  ForgotPasswordData,
  ResetPasswordData,
  SignUpResponse,
  ConfirmSignUpResponse,
  LoginResponse,
  RefreshTokenResponse,
  ForgotPasswordResponse,
  ResetPasswordResponse,
} from '../types';

export const authApi = {
  signUp: async (data: SignUpData): Promise<SignUpResponse> => {
    return apiClient.post<SignUpResponse>('/api/v1/auth/signup', data);
  },

  confirmSignUp: async (data: ConfirmSignUpData): Promise<ConfirmSignUpResponse> => {
    return apiClient.post<ConfirmSignUpResponse>('/api/v1/auth/confirm', data);
  },

  login: async (credentials: LoginCredentials): Promise<LoginResponse> => {
    return apiClient.post<LoginResponse>('/api/v1/auth/login', credentials);
  },

  refreshToken: async (refreshToken: string, email: string): Promise<RefreshTokenResponse> => {
    return apiClient.post<RefreshTokenResponse>('/api/v1/auth/refresh', {
      refresh_token: refreshToken,
      email,
    });
  },

  forgotPassword: async (data: ForgotPasswordData): Promise<ForgotPasswordResponse> => {
    return apiClient.post<ForgotPasswordResponse>('/api/v1/auth/forgot-password', data);
  },

  resetPassword: async (data: ResetPasswordData): Promise<ResetPasswordResponse> => {
    return apiClient.post<ResetPasswordResponse>('/api/v1/auth/reset-password', {
      email: data.email,
      code: data.code,
      new_password: data.newPassword,
    });
  },
};
