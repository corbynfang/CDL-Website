import { useState, useEffect } from 'react';
import axios, { AxiosError } from 'axios';

interface UseApiOptions {
  retries?: number;
  retryDelay?: number;
  enabled?: boolean;
}

interface UseApiResult<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

/**
 * Custom hook for API calls with automatic retry logic and error handling
 * 
 * @param url - API endpoint to fetch
 * @param options - Configuration options
 * @returns Object containing data, loading state, error, and refetch function
 * 
 * @example
 * const { data, loading, error } = useApi<Player>('/api/v1/players/1');
 */
export function useApi<T>(
  url: string,
  options: UseApiOptions = {}
): UseApiResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refetchTrigger, setRefetchTrigger] = useState(0);

  const { 
    retries = 3, 
    retryDelay = 1000,
    enabled = true 
  } = options;

  useEffect(() => {
    if (!enabled) {
      setLoading(false);
      return;
    }

    let isMounted = true;
    let retryCount = 0;
    const controller = new AbortController();

    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);

        const response = await axios.get<T>(url, {
          signal: controller.signal,
          timeout: 10000, // 10 second timeout
        });

        if (isMounted) {
          setData(response.data);
          setError(null);
        }
      } catch (err) {
        // Don't handle errors if request was cancelled
        if (axios.isCancel(err)) {
          return;
        }

        const axiosError = err as AxiosError;

        // Retry logic for network errors (not 404s)
        if (retryCount < retries && axiosError.response?.status !== 404) {
          retryCount++;
          console.log(`Retrying API call (${retryCount}/${retries})...`);
          setTimeout(fetchData, retryDelay * retryCount);
          return;
        }

        // Set appropriate error message
        if (isMounted) {
          if (axiosError.response?.status === 404) {
            setError('Not found');
          } else if (axiosError.response?.status === 429) {
            setError('Too many requests. Please try again later.');
          } else if (axiosError.code === 'ECONNABORTED') {
            setError('Request timeout. Please check your connection.');
          } else if (!axiosError.response) {
            setError('Network error. Please check your connection.');
          } else {
            setError('Failed to load data. Please try again.');
          }
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      isMounted = false;
      controller.abort();
    };
  }, [url, retries, retryDelay, enabled, refetchTrigger]);

  const refetch = () => {
    setRefetchTrigger((prev) => prev + 1);
  };

  return { data, loading, error, refetch };
}

