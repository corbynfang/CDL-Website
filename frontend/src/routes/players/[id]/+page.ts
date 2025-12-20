import { api, type Player, type PlayerKDResponse, type PlayerMatchesResponse } from '$lib/api';
import type { PageLoad } from './$types';

// Dynamic route params come from the URL
// /players/5 -> params.id = "5"
export const load: PageLoad = async ({ params }) => {
	const id = parseInt(params.id);
	
	if (isNaN(id)) {
		return {
			player: null,
			stats: null,
			matches: null,
			error: 'Invalid player ID'
		};
	}
	
	try {
		// Fetch all data in parallel (like Promise.all)
		const [player, stats, matches] = await Promise.all([
			api.getPlayer(id),
			api.getPlayerKD(id),
			api.getPlayerMatches(id)
		]);
		
		return {
			player,
			stats,
			matches,
			error: null
		};
	} catch (e) {
		console.error('Failed to fetch player data:', e);
		return {
			player: null,
			stats: null,
			matches: null,
			error: 'Player not found'
		};
	}
};

