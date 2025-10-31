/**
 * API Configuration
 * Centralizes API URL configuration and provides helper functions
 */

export const API_CONFIG = {
  // Base URL - automatically picks correct URL based on environment
  BASE_URL: import.meta.env.VITE_API_URL || 'https://cdlytics.me/api/v1',

  // Timeout settings
  TIMEOUT: 10000, // 10 seconds

  // Retry settings
  RETRY_ATTEMPTS: 3,
  RETRY_DELAY: 1000, // 1 second

  // Cache settings
  CACHE_TIME: 5 * 60 * 1000, // 5 minutes
  STALE_TIME: 2 * 60 * 1000, // 2 minutes
} as const;

/**
 * Build full API URL from endpoint
 */
export const getApiUrl = (endpoint: string): string => {
  // Remove leading slash if present
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint.slice(1) : endpoint;

  // Ensure base URL doesn't have trailing slash
  const baseUrl = API_CONFIG.BASE_URL.replace(/\/$/, '');

  return `${baseUrl}/${cleanEndpoint}`;
};

/**
 * Check if running in development mode
 */
export const isDevelopment = (): boolean => {
  return import.meta.env.DEV;
};

/**
 * Check if running in production mode
 */
export const isProduction = (): boolean => {
  return import.meta.env.PROD;
};

/**
 * Get environment name
 */
export const getEnvironment = (): 'development' | 'production' => {
  return isProduction() ? 'production' : 'development';
};

export default API_CONFIG;
