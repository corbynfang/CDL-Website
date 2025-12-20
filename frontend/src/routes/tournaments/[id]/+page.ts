import { api, type BracketResponse, type Tournament } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	const id = parseInt(params.id);
	
	if (isNaN(id)) {
		return { 
			tournament: null, 
			bracketData: null, 
			error: 'Invalid tournament ID' 
		};
	}
	
	try {
		const [tournament, bracketData] = await Promise.all([
			api.getTournament(id),
			api.getTournamentBracket(id)
		]);
		
		return { 
			tournament, 
			bracketData, 
			error: null 
		};
	} catch (e) {
		console.error('Failed to fetch tournament data:', e);
		return { 
			tournament: null, 
			bracketData: null, 
			error: 'Tournament not found' 
		};
	}
};

