import { writable, get } from 'svelte/store';
import { browser } from '$app/environment';

// =============================================================================
// THEME STORE
// =============================================================================

// Check for saved preference or system preference
function getInitialTheme(): 'light' | 'dark' {
	if (!browser) return 'light';
	
	const saved = localStorage.getItem('theme');
	if (saved === 'dark' || saved === 'light') return saved;
	
	// Check system preference
	if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
		return 'dark';
	}
	
	return 'light';
}

// Create the theme store
function createThemeStore() {
	const { subscribe, set, update } = writable<'light' | 'dark'>(getInitialTheme());
	
	return {
		subscribe,
		toggle: () => {
			update(current => {
				const newTheme = current === 'light' ? 'dark' : 'light';
				if (browser) {
					localStorage.setItem('theme', newTheme);
					document.documentElement.classList.toggle('dark', newTheme === 'dark');
				}
				return newTheme;
			});
		},
		set: (value: 'light' | 'dark') => {
			if (browser) {
				localStorage.setItem('theme', value);
				document.documentElement.classList.toggle('dark', value === 'dark');
			}
			set(value);
		},
		init: () => {
			if (browser) {
				const theme = getInitialTheme();
				document.documentElement.classList.toggle('dark', theme === 'dark');
				set(theme);
			}
		}
	};
}

export const theme = createThemeStore();

// =============================================================================
// GAME/SEASON STORE
// =============================================================================

export interface Season {
	id: number;
	name: string;
	game_title: string;
	start_date: string;
	end_date?: string;
	is_active: boolean;
}

// Available games with their metadata
export const GAMES = {
	'bo6': {
		id: 'bo6',
		name: 'Black Ops 6',
		shortName: 'BO6',
		year: '2024-2025',
		color: '#ff6b00'
	},
	'bo7': {
		id: 'bo7',
		name: 'Black Ops 7',
		shortName: 'BO7',
		year: '2025-2026',
		color: '#00d4ff'
	}
} as const;

export type GameId = keyof typeof GAMES;

function getInitialGame(): GameId {
	if (!browser) return 'bo6';
	const saved = localStorage.getItem('selectedGame');
	if (saved && (saved === 'bo6' || saved === 'bo7')) return saved;
	return 'bo6'; // Default to BO6 since it has data
}

function createGameStore() {
	const { subscribe, set, update } = writable<GameId>(getInitialGame());
	const seasonsStore = writable<Season[]>([]);
	const loadingStore = writable<boolean>(false);

	return {
		subscribe,
		seasons: { subscribe: seasonsStore.subscribe },
		loading: { subscribe: loadingStore.subscribe },
		
		select: (gameId: GameId) => {
			if (browser) {
				localStorage.setItem('selectedGame', gameId);
			}
			set(gameId);
		},
		
		init: () => {
			if (browser) {
				const game = getInitialGame();
				set(game);
			}
		},

		// Get current game info
		getGameInfo: () => {
			const currentGame = get({ subscribe });
			return GAMES[currentGame];
		},

		// Get season ID for API calls (maps game to season_id)
		getSeasonId: (): number | null => {
			const currentGame = get({ subscribe });
			// These IDs should match your database seasons
			// BO6 = season 1, BO7 = season 2 (adjust as needed)
			const seasonMap: Record<GameId, number> = {
				'bo6': 1,
				'bo7': 2
			};
			return seasonMap[currentGame] || null;
		},

		// Load seasons from API
		loadSeasons: async (apiBaseUrl: string) => {
			loadingStore.set(true);
			try {
				const response = await fetch(`${apiBaseUrl}/seasons`);
				if (response.ok) {
					const seasons = await response.json();
					seasonsStore.set(seasons);
				}
			} catch (e) {
				console.error('Failed to load seasons:', e);
			} finally {
				loadingStore.set(false);
			}
		}
	};
}

export const gameStore = createGameStore();