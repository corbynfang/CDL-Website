import { api, type Tournament } from '$lib/api';
import { browser } from '$app/environment';
import type { PageLoad } from './$types';

// Get the season ID from localStorage (matches the game store logic)
function getSeasonId(): number {
	if (!browser) return 1; // Default to BO6 on server
	const savedGame = localStorage.getItem('selectedGame');
	const seasonMap: Record<string, number> = { 'bo6': 1, 'bo7': 2 };
	return seasonMap[savedGame || 'bo6'] || 1;
}

export const load: PageLoad = async ({ depends }) => {
	// Tell SvelteKit this load depends on 'game' so invalidateAll() refreshes it
	depends('game');
	
	const seasonId = getSeasonId();
	
	try {
		const tournaments = await api.getTournaments(seasonId);
		return { tournaments, error: null, seasonId };
	} catch (e) {
		console.error('Failed to fetch tournaments:', e);
		return { tournaments: [] as Tournament[], error: 'Failed to load tournaments', seasonId };
	}
};

