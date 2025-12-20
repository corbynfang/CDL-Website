import { api, type Tournament } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const tournaments = await api.getTournaments();
		return { tournaments, error: null };
	} catch (e) {
		console.error('Failed to fetch tournaments:', e);
		return { tournaments: [] as Tournament[], error: 'Failed to load tournaments' };
	}
};

