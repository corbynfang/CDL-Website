import { api, type StatsResponse } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const statsData = await api.getAllKDStats();
		return { statsData, error: null };
	} catch (e) {
		console.error('Failed to fetch stats:', e);
		return { statsData: null, error: 'Failed to load statistics' };
	}
};

