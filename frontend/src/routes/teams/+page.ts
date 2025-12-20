import { api, type Team } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const teams = await api.getTeams();
		return { teams, error: null };
	} catch (e) {
		console.error('Failed to fetch teams:', e);
		return { teams: [] as Team[], error: 'Failed to load teams' };
	}
};

