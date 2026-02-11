import { computed } from 'vue';
import { useStorage } from '@vueuse/core';
import { API_BASE_URL } from '../config';

export interface User {
  id: number;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  refreshToken: string;
  user: User;
}

// Store authentication state
const token = useStorage<string | null>('auth_token', null);
const refreshToken = useStorage<string | null>('auth_refresh_token', null);
const user = useStorage<User | null>('auth_user', null);

// Computed for authentication status
export const isAuthenticated = computed(() => !!token.value && !!user.value);
export const currentUser = computed(() => user.value);

// Get authentication token
export function getAuthToken(): string | null {
  return token.value;
}

// API call helper with authentication
export async function apiCall(endpoint: string, options?: RequestInit, isRetry: boolean = false) {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options?.headers as Record<string, string>,
  };

  // Add authentication token if available
  if (token.value) {
    headers['Authorization'] = `Bearer ${token.value}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  // Handle 401 Unauthorized - try to refresh token (only if not already retrying)
  if (response.status === 401 && refreshToken.value && !isRetry) {
    const refreshed = await tryRefreshToken();
    if (refreshed) {
      // Retry the original request with new token (mark as retry to prevent infinite loop)
      headers['Authorization'] = `Bearer ${token.value}`;
      const retryResponse = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
      });
      
      if (!retryResponse.ok) {
        // If retry also fails, logout
        if (retryResponse.status === 401) {
          logout();
          throw new Error('Session expired. Please login again.');
        }
        throw new Error(`API call failed: ${retryResponse.statusText}`);
      }
      
      return retryResponse.json();
    } else {
      // Refresh failed, logout user
      logout();
      throw new Error('Session expired. Please login again.');
    }
  }

  if (!response.ok) {
    throw new Error(`API call failed: ${response.statusText}`);
  }

  return response.json();
}

// Register a new user
export async function register(email: string, password: string): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/api/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Registration failed');
  }

  const data: AuthResponse = await response.json();
  setAuthData(data);
  return data;
}

// Login user
export async function login(email: string, password: string): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Login failed');
  }

  const data: AuthResponse = await response.json();
  setAuthData(data);
  return data;
}

// Logout user
export function logout() {
  token.value = null;
  refreshToken.value = null;
  user.value = null;
}

// Try to refresh the authentication token
async function tryRefreshToken(): Promise<boolean> {
  if (!refreshToken.value) return false;

  try {
    const response = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refreshToken: refreshToken.value }),
    });

    if (!response.ok) return false;

    const data = await response.json();
    token.value = data.token;
    refreshToken.value = data.refreshToken;
    return true;
  } catch {
    return false;
  }
}

// Set authentication data
function setAuthData(data: AuthResponse) {
  token.value = data.token;
  refreshToken.value = data.refreshToken;
  user.value = data.user;
}
