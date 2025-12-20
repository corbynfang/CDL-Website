import { api, type TransfersResponse } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const response = await api.getTransfers();
		return { transfers: response.transfers, error: null };
	} catch (e) {
		console.error('Failed to fetch transfers:', e);
		return { transfers: [], error: 'Failed to load transfers' };
	}
};

