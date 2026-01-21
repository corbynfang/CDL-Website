<script lang="ts">
	import type { PlayerStats } from '$lib/api';
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
	const statsData = $derived(data.statsData);
	const error = $derived(data.error as string | null);
	
	// Get players array and sort by K/D
	const players = $derived(
		(statsData?.players || [])
			.slice()
			.sort((a: PlayerStats, b: PlayerStats) => b.season_kd - a.season_kd)
	);
	
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
</script>

<svelte:head>
	<title>K/D Leaderboard - {GAMES[currentGame].shortName} - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<div class="flex items-center gap-4 mb-8">
			<h1 class="text-4xl font-bold">K/D LEADERBOARD</h1>
			<span 
				class="px-3 py-1 text-sm font-semibold rounded"
				style="background-color: {GAMES[currentGame].color}20; color: {GAMES[currentGame].color}; border: 1px solid {GAMES[currentGame].color}"
			>
				{GAMES[currentGame].shortName}
			</span>
		</div>
		
		{#if error}
			<p class="text-red-500">{error}</p>
		{:else if players.length === 0}
			<!-- Empty state for BO7 or games with no stats yet -->
			<div class="text-center py-16 bg-theme-secondary rounded-lg">
				<div class="text-6xl mb-4">📊</div>
				<h2 class="text-2xl font-bold mb-2">No Stats Available Yet</h2>
				<p class="text-theme-secondary max-w-md mx-auto">
					Statistics for {GAMES[currentGame].name} ({GAMES[currentGame].year}) 
					will be available once the season begins and matches are played.
				</p>
				{#if currentGame === 'bo7'}
					<div class="mt-6 p-4 bg-blue-500/10 border border-blue-500/30 rounded-lg inline-block">
						<p class="text-blue-400 text-sm">
							📈 Stats will populate after the first Black Ops 7 CDL tournament
						</p>
					</div>
				{/if}
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="border-b border-theme">
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider w-16">
								Rank
							</th>
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Player
							</th>
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Team
							</th>
							<th class="text-right py-4 text-theme-secondary text-sm uppercase tracking-wider">
								K/D
							</th>
							<th class="text-right py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Kills
							</th>
							<th class="text-right py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Deaths
							</th>
							<th class="text-right py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Assists
							</th>
						</tr>
					</thead>
					<tbody>
						{#each players as player, index (player.player_id)}
							<tr class="border-b border-theme hover:bg-theme-secondary/50 transition-colors">
								<td class="py-4 text-theme-secondary">{index + 1}</td>
								<td class="py-4">
									<a 
										href="/players/{player.player_id}" 
										class="font-semibold hover:opacity-70 transition-opacity"
									>
										{player.gamertag}
									</a>
								</td>
								<td class="py-4 text-theme-secondary">
									{player.team_abbr || '—'}
								</td>
								<td class="py-4 text-right font-bold">
									{player.season_kd?.toFixed(2) || '0.00'}
								</td>
								<td class="py-4 text-right text-theme-secondary">
									{player.season_kills?.toLocaleString() || '0'}
								</td>
								<td class="py-4 text-right text-theme-secondary">
									{player.season_deaths?.toLocaleString() || '0'}
								</td>
								<td class="py-4 text-right text-theme-secondary">
									{player.season_assists?.toLocaleString() || '0'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			
			<p class="mt-8 text-theme-secondary text-sm">
				Total Players: {players.length}
			</p>
		{/if}
	</div>
</div>

