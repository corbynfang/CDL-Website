<script lang="ts">
	import type { Tournament } from '$lib/api';
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	
	let { data } = $props();
	
	// Game configuration
	const GAMES = {
		'bo6': { id: 'bo6', name: 'Black Ops 6', shortName: 'BO6', year: '2024-2025', color: '#ff6b00' },
		'bo7': { id: 'bo7', name: 'Black Ops 7', shortName: 'BO7', year: '2025-2026', color: '#00d4ff' }
	} as const;
	
	type GameId = 'bo6' | 'bo7';
	
	// Use $derived for reactive access to data properties
	const tournaments = $derived(data.tournaments as Tournament[]);
	const error = $derived(data.error as string | null);
	
	// Track current game
	let currentGame = $state<GameId>('bo6');
	
	onMount(() => {
		if (browser) {
			const saved = localStorage.getItem('selectedGame');
			if (saved === 'bo6' || saved === 'bo7') {
				currentGame = saved;
			}
		}
	});
	
	function formatDate(dateString: string): string {
		try {
			return new Date(dateString).toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'short',
				day: 'numeric'
			});
		} catch {
			return dateString;
		}
	}
</script>

<svelte:head>
	<title>Tournaments - {GAMES[currentGame].shortName} - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<div class="flex items-center gap-4 mb-8">
			<h1 class="text-4xl font-bold">TOURNAMENTS</h1>
			<span 
				class="px-3 py-1 text-sm font-semibold rounded"
				style="background-color: {GAMES[currentGame].color}20; color: {GAMES[currentGame].color}; border: 1px solid {GAMES[currentGame].color}"
			>
				{GAMES[currentGame].shortName}
			</span>
		</div>
		
		{#if error}
			<p class="text-red-500">{error}</p>
		{:else if tournaments.length === 0}
			<!-- Empty state for BO7 or games with no tournaments yet -->
			<div class="text-center py-16 bg-theme-secondary rounded-lg">
				<div class="text-6xl mb-4">🏆</div>
				<h2 class="text-2xl font-bold mb-2">No Tournaments Yet</h2>
				<p class="text-theme-secondary max-w-md mx-auto">
					Tournaments for {GAMES[currentGame].name} ({GAMES[currentGame].year}) 
					haven't been scheduled yet. Check back soon!
				</p>
				{#if currentGame === 'bo7'}
					<div class="mt-6 p-4 bg-blue-500/10 border border-blue-500/30 rounded-lg inline-block">
						<p class="text-blue-400 text-sm">
							🗓️ Black Ops 7 CDL tournament schedule coming soon
						</p>
					</div>
				{/if}
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each tournaments as tournament (tournament.id)}
					<a
						href="/tournaments/{tournament.id}"
						class="p-6 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200 group"
					>
						<div class="mb-4">
							<h2 class="text-xl font-bold group-hover:text-green-500 transition-colors">{tournament.name}</h2>
							{#if tournament.tournament_type}
								<span class="text-xs uppercase tracking-wider text-theme-secondary">{tournament.tournament_type}</span>
							{/if}
						</div>
						
						<div class="space-y-2 text-sm text-theme-secondary">
							{#if tournament.location}
								<p>📍 {tournament.location}</p>
							{/if}
							<p>📅 {formatDate(tournament.start_date)}</p>
							{#if tournament.prize_pool}
								<p>💰 ${tournament.prize_pool.toLocaleString()}</p>
							{/if}
						</div>
						
						<div class="mt-4 pt-4 border-t border-theme">
							<span class="text-sm text-green-500 font-medium">View Bracket →</span>
						</div>
					</a>
				{/each}
			</div>
			
			<p class="mt-8 text-theme-secondary text-sm">
				Total Tournaments: {tournaments.length}
			</p>
		{/if}
	</div>
</div>

