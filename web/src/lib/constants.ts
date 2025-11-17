export const APP_NAME = 'AWS Go Server';

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  PROFILE: '/profile',
  SETTINGS: '/settings',
} as const;

export const QUERY_KEYS = {
  USER: 'user',
  ITEMS: 'items',
  DASHBOARD: 'dashboard',
  PROFILE: 'profile',
} as const;

export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'access_token',
  ID_TOKEN: 'id_token',
  REFRESH_TOKEN: 'refresh_token',
  USER_EMAIL: 'user_email',
} as const;
