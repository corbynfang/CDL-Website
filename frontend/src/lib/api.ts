// API configuration and helper functions
// This is the SvelteKit equivalent of your React api.ts

// API URL - uses production by default, override with VITE_API_URL for local dev
const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://cdlytics.me/api/v1';

export { API_BASE_URL };

// Types matching your Go backend models
export interface Player {
	id: number;
	gamertag: string;
	first_name?: string;
	last_name?: string;
	country?: string;
	birthdate?: string;
	role?: string;
	is_active: boolean;
	liquipedia_url?: string;
	twitter_handle?: string;
	avatar_url?: string;
	created_at: string;
	updated_at: string;
}

export interface Team {
	id: number;
	name: string;
	abbreviation: string;
	city?: string;
	logo_url?: string;
	primary_color?: string;
	secondary_color?: string;
	founded_date?: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface PlayerStats {
	player_id: number;
	gamertag: string;
	team_abbr: string;
	season_kd: number;
	season_kills: number;
	season_deaths: number;
	season_assists: number;
}

export interface PlayerTransfer {
	id: number;
	player_id: number;
	from_team_id?: number;
	to_team_id?: number;
	transfer_date: string;
	transfer_type: string;
	role: string;
	season: string;
	description: string;
	created_at: string;
	player?: Player;
	from_team?: Team;
	to_team?: Team;
}

export interface TransfersResponse {
	transfers: PlayerTransfer[];
	count: number;
	timestamp?: number;
}

export interface StatsResponse {
	players: PlayerStats[];
	count: number;
	timestamp: number;
}

export interface PlayerKDResponse {
	player_id: number;
	gamertag: string;
	avatar_url?: string;
	total_kills: number;
	total_deaths: number;
	total_assists: number;
	avg_kd: number;
	avg_kda: number;
	hp_kd_ratio: number;
	snd_kd_ratio: number;
	control_kd_ratio: number;
	tournament_stats: any[];
}

export interface PlayerMatchesResponse {
	player_id: number;
	events: any[];
	total: number;
}

export interface Tournament {
	id: number;
	name: string;
	tournament_type?: string;
	start_date: string;
	end_date?: string;
	prize_pool?: number;
	location?: string;
	tournament_format?: string;
}

export interface BracketMatch {
	id: number;
	team1_id: number;
	team2_id: number;
	team1_name: string;
	team1_abbr: string;
	team1_logo: string;
	team2_name: string;
	team2_abbr: string;
	team2_logo: string;
	team1_score: number;
	team2_score: number;
	winner_id: number | null;
	bracket_position: number;
	match_date: string;
}

export interface BracketData {
	winners_r1: BracketMatch[];
	winners_r2: BracketMatch[];
	winners_finals: BracketMatch[];
	elim_r1: BracketMatch[];
	elim_r2: BracketMatch[];
	elim_r3: BracketMatch[];
	elim_finals: BracketMatch[];
	grand_finals: BracketMatch[];
}

export interface BracketResponse {
	tournament_id: number;
	tournament_name: string;
	bracket: BracketData;
	total_matches: number;
}

// Season type
export interface Season {
	id: number;
	name: string;
	game_title: string;
	start_date: string;
	end_date?: string;
	is_active: boolean;
}

// Generic fetch wrapper with cache busting
async function fetchAPI<T>(endpoint: string, params?: Record<string, string | number | undefined>): Promise<T> {
	const timestamp = Date.now();
	const random = Math.random().toString(36).substring(7);
	
	// Build query string with optional parameters
	const searchParams = new URLSearchParams();
	searchParams.set('_t', timestamp.toString());
	searchParams.set('_r', random);
	
	if (params) {
		Object.entries(params).forEach(([key, value]) => {
			if (value !== undefined && value !== null) {
				searchParams.set(key, value.toString());
			}
		});
	}
	
	const url = `${API_BASE_URL}${endpoint}?${searchParams.toString()}`;
	
	const response = await fetch(url, {
		headers: {
			'Content-Type': 'application/json',
			'Cache-Control': 'no-cache, no-store, must-revalidate',
		},
	});
	
	if (!response.ok) {
		throw new Error(`API Error: ${response.status} ${response.statusText}`);
	}
	
	return response.json();
}

// API functions
export const api = {
	// Seasons
	getSeasons: () => fetchAPI<Season[]>('/seasons'),
	getSeason: (id: number) => fetchAPI<Season>(`/seasons/${id}`),
	getActiveSeason: () => fetchAPI<Season>('/seasons/active'),

	// Players
	getPlayers: () => fetchAPI<Player[]>('/players'),
	getPlayer: (id: number) => fetchAPI<Player>(`/players/${id}`),
	getPlayerKD: (id: number) => fetchAPI<PlayerKDResponse>(`/players/${id}/kd`),
	getPlayerMatches: (id: number) => fetchAPI<PlayerMatchesResponse>(`/players/${id}/matches`),
	
	// Teams (supports season filtering)
	getTeams: (seasonId?: number) => fetchAPI<Team[]>('/teams', { season_id: seasonId }),
	getTeam: (id: number) => fetchAPI<Team>(`/teams/${id}`),
	getTeamPlayers: (id: number, seasonId?: number) => fetchAPI<Player[]>(`/teams/${id}/players`, { season_id: seasonId }),
	
	// Tournaments (supports season filtering)
	getTournaments: (seasonId?: number) => fetchAPI<Tournament[]>('/tournaments', { season_id: seasonId }),
	getTournament: (id: number) => fetchAPI<Tournament>(`/tournaments/${id}`),
	getTournamentBracket: (id: number) => fetchAPI<BracketResponse>(`/tournaments/${id}/bracket`),
	
	// Stats (supports season filtering)
	getAllKDStats: (seasonId?: number) => fetchAPI<StatsResponse>('/stats/all-kd-by-tournament', { season_id: seasonId }),
	
	// Transfers
	getTransfers: () => fetchAPI<TransfersResponse>('/transfers'),
};

