import type { CognitoUser, User } from '../types';

/**
 * Decode a JWT token (without verification - server validates)
 * This is safe for client-side use as we're just extracting claims
 */
export const decodeJWT = <T = any>(token: string): T | null => {
  try {
    const parts = token.split('.');
    if (parts.length !== 3) {
      return null;
    }

    const payload = parts[1];
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
    return JSON.parse(decoded);
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
};

/**
 * Extract user information from Cognito ID token
 */
export const getUserFromIdToken = (idToken: string): User | null => {
  const cognitoUser = decodeJWT<CognitoUser>(idToken);
  if (!cognitoUser) {
    return null;
  }

  return {
    id: cognitoUser.sub,
    email: cognitoUser.email,
    username: cognitoUser['cognito:username'] || cognitoUser.email.split('@')[0],
    name: cognitoUser.name,
    firstName: cognitoUser.given_name,
    lastName: cognitoUser.family_name,
    emailVerified: cognitoUser.email_verified,
    groups: cognitoUser['cognito:groups'],
  };
};

/**
 * Check if a JWT token is expired
 */
export const isTokenExpired = (token: string): boolean => {
  const decoded = decodeJWT<{ exp?: number }>(token);
  if (!decoded || !decoded.exp) {
    return true;
  }

  const currentTime = Math.floor(Date.now() / 1000);
  return decoded.exp < currentTime;
};

/**
 * Get token expiration time in seconds
 */
export const getTokenExpiration = (token: string): number | null => {
  const decoded = decodeJWT<{ exp?: number }>(token);
  return decoded?.exp || null;
};
