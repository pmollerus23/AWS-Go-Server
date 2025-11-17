import type { ApiResponse, ApiError } from '../types';
import { storage, STORAGE_KEYS } from '../lib';

interface RequestConfig extends RequestInit {
  params?: Record<string, string | number | boolean>;
}

class ApiClient {
  private baseURL: string;
  private defaultHeaders: HeadersInit;

  constructor(baseURL: string = import.meta.env.VITE_API_URL || '') {
    this.baseURL = baseURL;
    this.defaultHeaders = {
      'Content-Type': 'application/json',
    };
  }

  private getAuthToken(): string | null {
    return storage.get<string>(STORAGE_KEYS.ACCESS_TOKEN);
  }

  private buildURL(endpoint: string, params?: Record<string, string | number | boolean>): string {
    // If baseURL is empty, use relative URLs (same origin)
    if (!this.baseURL) {
      if (params) {
        const searchParams = new URLSearchParams();
        Object.entries(params).forEach(([key, value]) => {
          searchParams.append(key, String(value));
        });
        const queryString = searchParams.toString();
        return queryString ? `${endpoint}?${queryString}` : endpoint;
      }
      return endpoint;
    }

    // If baseURL is provided, construct full URL
    const url = new URL(endpoint, this.baseURL);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        url.searchParams.append(key, String(value));
      });
    }

    return url.toString();
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    const contentType = response.headers.get('content-type');
    const isJson = contentType?.includes('application/json');

    if (!response.ok) {
      const errorData = isJson ? await response.json() : { message: response.statusText };

      const error: ApiError = {
        message: errorData.message || 'An error occurred',
        status: response.status,
        code: errorData.code,
        details: errorData.details,
      };

      throw error;
    }

    if (isJson) {
      const data = await response.json();
      // Handle wrapped responses
      if (data && typeof data === 'object' && 'data' in data) {
        return (data as ApiResponse<T>).data;
      }
      return data;
    }

    return {} as T;
  }

  private async request<T>(
    endpoint: string,
    config: RequestConfig = {}
  ): Promise<T> {
    const { params, headers, ...restConfig } = config;

    const url = this.buildURL(endpoint, params);
    const token = this.getAuthToken();

    const requestHeaders: HeadersInit = {
      ...this.defaultHeaders,
      ...headers,
    };

    if (token) {
      (requestHeaders as Record<string, string>)['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...restConfig,
      headers: requestHeaders,
    });

    return this.handleResponse<T>(response);
  }

  async get<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, { ...config, method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: unknown, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async patch<T>(endpoint: string, data?: unknown, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, { ...config, method: 'DELETE' });
  }

  setBaseURL(url: string): void {
    this.baseURL = url;
  }

  setDefaultHeader(key: string, value: string): void {
    (this.defaultHeaders as Record<string, string>)[key] = value;
  }
}

export const apiClient = new ApiClient();
