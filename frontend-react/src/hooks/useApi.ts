import { useState, useEffect } from "react";
import axios from "axios";
import type { AxiosError } from "axios";

interface UseApiOptions {
	retries?: number;
	retryDelay?: number;
	enabled?: boolean;
	cacheTtl?: number;
}

interface CacheEntry<T> {
	data: T;
	expiresAt: number;
}

const apiCache = new Map<string, CacheEntry<unknown>>();

interface UseApiResult<T> {
	data: T | null;
	loading: boolean;
	error: string | null;
	refetch: () => void;
}

export function useApi<T>(
	url: string,
	options: UseApiOptions = {},
): UseApiResult<T> {
	const [data, setData] = useState<T | null>(null);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [refetchTrigger, setRefetchTrigger] = useState(0);

	const { retries = 3, retryDelay = 1000, enabled = true, cacheTtl } = options;

	useEffect(() => {
		if (!enabled) {
			setLoading(false);
			return;
		}

		const skipCache = refetchTrigger > 0;
		if (!skipCache && cacheTtl !== undefined) {
			const entry = apiCache.get(url);
			if (entry && Date.now() < entry.expiresAt) {
				setData(entry.data as T);
				setLoading(false);
				return;
			}
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
					if (cacheTtl !== undefined) {
						apiCache.set(url, {
							data: response.data,
							expiresAt: Date.now() + cacheTtl,
						});
					}
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
						setError("Not found");
					} else if (axiosError.response?.status === 429) {
						setError("Too many requests. Please try again later.");
					} else if (axiosError.code === "ECONNABORTED") {
						setError("Request timeout. Please check your connection.");
					} else if (!axiosError.response) {
						setError("Network error. Please check your connection.");
					} else {
						setError("Failed to load data. Please try again.");
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
	}, [url, retries, retryDelay, enabled, cacheTtl, refetchTrigger]);

	const refetch = () => {
		setRefetchTrigger((prev) => prev + 1);
	};

	return { data, loading, error, refetch };
}
