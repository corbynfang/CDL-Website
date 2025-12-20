import { api, type Player } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const players = await api.getPlayers();
		return {
			players,
			error: null
		};
	} catch (e) {
		console.error('Failed to fetch players:', e);
		return {
			players: [] as Player[],
			error: 'Failed to load players'
		};
	}
};

