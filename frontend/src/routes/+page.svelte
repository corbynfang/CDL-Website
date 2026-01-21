<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { invalidateAll } from '$app/navigation';
	
	// Game configuration
	const GAMES = {
		'bo6': { id: 'bo6', name: 'Black Ops 6', shortName: 'BO6', year: '2024-2025', color: '#ff6b00' },
		'bo7': { id: 'bo7', name: 'Black Ops 7', shortName: 'BO7', year: '2025-2026', color: '#00d4ff' }
	} as const;
	
	type GameId = 'bo6' | 'bo7';
	
	// Track current game
	let currentGame = $state<GameId>('bo6');
	
	onMount(() => {
		if (browser) {
			const saved = localStorage.getItem('selectedGame');
			if (saved === 'bo6' || saved === 'bo7') {
				currentGame = saved;
			}
			
			// Listen for storage changes (when user switches game in another tab/component)
			const handleStorage = (e: StorageEvent) => {
				if (e.key === 'selectedGame' && (e.newValue === 'bo6' || e.newValue === 'bo7')) {
					currentGame = e.newValue;
				}
			};
			window.addEventListener('storage', handleStorage);
			return () => window.removeEventListener('storage', handleStorage);
		}
	});
	
	function switchToBO6() {
		if (browser) {
			localStorage.setItem('selectedGame', 'bo6');
			currentGame = 'bo6';
			invalidateAll();
		}
	}
</script>

<svelte:head>
	<title>CDLYTICS - {GAMES[currentGame].name} Statistics</title>
	<meta name="description" content="Call of Duty League Statistics & Analytics for {GAMES[currentGame].name}" />
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<!-- Hero Section -->
		<div class="text-center mb-16">
			<h1 class="text-5xl font-bold mb-4">CDLYTICS</h1>
			<p class="text-xl text-theme-secondary mb-6">
				Call of Duty League Statistics & Analytics
			</p>
			
			<!-- Current Game Badge -->
			<div class="inline-flex items-center gap-3 px-6 py-3 bg-theme-secondary rounded-lg border border-theme">
				<div 
					class="w-3 h-3 rounded-full animate-pulse"
					style="background-color: {GAMES[currentGame].color}"
				></div>
				<span class="font-semibold">{GAMES[currentGame].name}</span>
				<span class="text-theme-secondary text-sm">{GAMES[currentGame].year}</span>
			</div>
		</div>
		
		<!-- Navigation Cards -->
		<div class="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
			<a 
				href="/players" 
				class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 group"
			>
				<div class="text-4xl mb-4">👤</div>
				<h2 class="text-2xl font-bold mb-2 group-hover:text-blue-400 transition-colors">PLAYERS</h2>
				<p class="text-theme-secondary">
					View player statistics and performance data
				</p>
			</a>
			
			<a 
				href="/teams" 
				class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 group"
			>
				<div class="text-4xl mb-4">🎮</div>
				<h2 class="text-2xl font-bold mb-2 group-hover:text-green-400 transition-colors">TEAMS</h2>
				<p class="text-theme-secondary">
					Browse {GAMES[currentGame].shortName} teams and rosters
				</p>
			</a>
			
			<a 
				href="/stats" 
				class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 group"
			>
				<div class="text-4xl mb-4">📊</div>
				<h2 class="text-2xl font-bold mb-2 group-hover:text-orange-400 transition-colors">STATS</h2>
				<p class="text-theme-secondary">
					K/D leaderboards and rankings
				</p>
			</a>
		</div>
		
		<!-- Quick Stats Preview for BO7 coming soon (After Major 1 Grand Finals) -->
		{#if currentGame === 'bo7'}
			<div class="mt-16 max-w-2xl mx-auto text-center">
				<div class="p-8 bg-gradient-to-br from-blue-500/10 to-cyan-500/10 border border-blue-500/30 rounded-lg">
					<div class="text-5xl mb-4">🚀</div>
					<h3 class="text-2xl font-bold mb-2">Black Ops 7 Season Coming Soon (After Major 1 Grand Finals)</h3>
					<p class="text-theme-secondary mb-4">
						The new CDL season is about to begin! Teams and stats will be updated 
						as the season progresses. (After Major 1 Grand Finals, Hopefully by then)
					</p>
					<div class="flex justify-center gap-4">
						<button 
							onclick={switchToBO6}
							class="px-4 py-2 bg-theme-primary hover:bg-theme-secondary border border-theme rounded transition-colors"
						>
							View BO6 Stats →
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
