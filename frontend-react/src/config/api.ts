export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_URL || "https://cdlytics.com/api/v1",

  TIMEOUT: 10000,

  RETRY_ATTEMPTS: 3,
  RETRY_DELAY: 1000,

  CACHE_TIME: 5 * 60 * 1000,
  STALE_TIME: 2 * 60 * 1000,
} as const;

export const getApiUrl = (endpoint: string): string => {
  const cleanEndpoint = endpoint.startsWith("/") ? endpoint.slice(1) : endpoint;
  const baseUrl = API_CONFIG.BASE_URL.replace(/\/$/, "");

  return `${baseUrl}/${cleanEndpoint}`;
};

export const isDevelopment = (): boolean => {
  return import.meta.env.DEV;
};

export const isProduction = (): boolean => {
  return import.meta.env.PROD;
};

export const getEnvironment = (): "development" | "production" => {
  return isProduction() ? "production" : "development";
};

export default API_CONFIG;
