<script lang="ts">
	import '../app.css';
	import { theme } from '$lib/stores/theme';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
	import { writable } from 'svelte/store';
	import { browser } from '$app/environment';
	
	let { children } = $props();
	
	// Game configuration - defined locally to avoid import issues
	const GAMES = {
		'bo6': { id: 'bo6', name: 'Black Ops 6', shortName: 'BO6', year: '2024-2025', color: '#ff6b00' },
		'bo7': { id: 'bo7', name: 'Black Ops 7', shortName: 'BO7', year: '2025-2026', color: '#00d4ff' }
	} as const;
	
	type GameId = 'bo6' | 'bo7';
	
	// Local game store
	function getInitialGame(): GameId {
		if (!browser) return 'bo6';
		const saved = localStorage.getItem('selectedGame');
		if (saved === 'bo6' || saved === 'bo7') return saved;
		return 'bo6';
	}
	
	const gameStore = writable<GameId>(getInitialGame());
	
	// Reactive subscription to game store
	let currentGame = $state<GameId>('bo6');
	let showGameMenu = $state(false);
	
	onMount(() => {
		theme.init();
		if (browser) {
			const saved = getInitialGame();
			gameStore.set(saved);
		}
		
		// Subscribe to game store changes
		const unsubscribe = gameStore.subscribe((value: GameId) => {
			currentGame = value;
		});
		
		return unsubscribe;
	});
	
	function selectGame(gameId: GameId) {
		if (browser) {
			localStorage.setItem('selectedGame', gameId);
		}
		gameStore.set(gameId);
		showGameMenu = false;
		// Refresh data when game changes
		invalidateAll();
	}
	
	function toggleGameMenu() {
		showGameMenu = !showGameMenu;
	}
	
	// Close menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.game-selector')) {
			showGameMenu = false;
		}
	}
</script>

<svelte:window onclick={handleClickOutside} />

<div class="min-h-screen bg-theme-primary text-theme-primary transition-colors duration-300">
	<!-- Game Selector Bar -->
	<div class="bg-theme-secondary border-b border-theme">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-10">
				<div class="flex items-center gap-2">
					<span class="text-xs uppercase tracking-wider text-theme-secondary">Game:</span>
					<div class="game-selector relative">
						<button 
							onclick={toggleGameMenu}
							class="flex items-center gap-2 px-3 py-1.5 text-sm font-semibold rounded transition-all hover:bg-theme-primary"
							style="border-left: 3px solid {GAMES[currentGame].color}"
						>
							<span>{GAMES[currentGame].shortName}</span>
							<svg class="w-4 h-4 transition-transform {showGameMenu ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
							</svg>
						</button>
						
						{#if showGameMenu}
							<div class="absolute top-full left-0 mt-1 w-56 bg-theme-secondary border border-theme rounded shadow-xl z-50">
								{#each Object.entries(GAMES) as [id, game]}
									<button
										onclick={() => selectGame(id as GameId)}
										class="w-full flex items-center gap-3 px-4 py-3 text-left hover:bg-theme-primary transition-colors {currentGame === id ? 'bg-theme-primary' : ''}"
									>
										<div 
											class="w-1 h-8 rounded"
											style="background-color: {game.color}"
										></div>
										<div>
											<div class="font-semibold text-sm">{game.name}</div>
											<div class="text-xs text-theme-secondary">{game.year}</div>
										</div>
										{#if currentGame === id}
											<svg class="w-4 h-4 ml-auto text-green-500" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
											</svg>
										{/if}
									</button>
								{/each}
							</div>
						{/if}
					</div>
				</div>
				
				<div class="text-xs text-theme-secondary">
					{#if currentGame === 'bo7'}
						<span class="px-2 py-0.5 bg-blue-500/20 text-blue-400 rounded">Season Starting Soon</span>
					{:else}
						<span>{GAMES[currentGame].year} Season</span>
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Header -->
	<header class="border-b border-theme">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-16">
				<a 
					href="/" 
					class="text-xl font-bold tracking-tight hover:opacity-80 transition-opacity"
				>
					CDLYTICS
				</a>
				
				<nav class="flex items-center space-x-8">
					<a 
						href="/players" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Players
					</a>
					<a 
						href="/teams" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Teams
					</a>
					<a 
						href="/tournaments" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Brackets
					</a>
					<a 
						href="/stats" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Stats
					</a>
					<a 
						href="/transfers" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Transfers
					</a>
					
					<!-- -->
					<ThemeToggle />
				</nav>
			</div>
		</div>
	</header>
	
	<!-- Main Content - this is where page content goes (like React's <Outlet />) -->
	<main>
		{@render children()}
	</main>
	
	<!-- Footer -->
	<footer class="border-t border-theme mt-16">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<p class="text-center text-theme-secondary text-sm">CDLytics © 2025</p>
		</div>
	</footer>
</div>
