import { api, type Team, type Player } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const id = parseInt(params.id);
	
	if (isNaN(id)) {
		return { team: null, players: [], error: 'Invalid team ID' };
	}
	
	try {
		const [team, players] = await Promise.all([
			api.getTeam(id),
			api.getTeamPlayers(id)
		]);
		
		return { team, players, error: null };
	} catch (e) {
		console.error('Failed to fetch team data:', e);
		return { team: null, players: [] as Player[], error: 'Team not found' };
	}
};

