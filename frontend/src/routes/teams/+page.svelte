<script lang="ts">
	import type { Team } from '$lib/api';
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
	const teams = $derived(data.teams as Team[]);
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
</script>

<svelte:head>
	<title>Teams - {GAMES[currentGame].shortName} - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<div class="flex items-center gap-4 mb-8">
			<h1 class="text-4xl font-bold">TEAMS</h1>
			<span 
				class="px-3 py-1 text-sm font-semibold rounded"
				style="background-color: {GAMES[currentGame].color}20; color: {GAMES[currentGame].color}; border: 1px solid {GAMES[currentGame].color}"
			>
				{GAMES[currentGame].shortName}
			</span>
		</div>
		
		{#if error}
			<p class="text-red-500">{error}</p>
		{:else if teams.length === 0}
			<!-- Empty state for BO7 or games with no teams yet -->
			<div class="text-center py-16 bg-theme-secondary rounded-lg">
				<div class="text-6xl mb-4">🎮</div>
				<h2 class="text-2xl font-bold mb-2">{GAMES[currentGame].name} Season Starting Soon</h2>
				<p class="text-theme-secondary max-w-md mx-auto">
					Teams for the {GAMES[currentGame].year} season haven't been announced yet. 
					Check back soon for roster updates!
				</p>
				{#if currentGame === 'bo7'}
					<div class="mt-6 p-4 bg-blue-500/10 border border-blue-500/30 rounded-lg inline-block">
						<p class="text-blue-400 text-sm">
							🔔 The Black Ops 7 CDL season is expected to start in late 2025
						</p>
					</div>
				{/if}
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each teams as team (team.id)}
					<a
						href="/teams/{team.id}"
						class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200"
					>
						<div class="flex items-center space-x-4 mb-4">
							{#if team.logo_url}
								<img
									src={team.logo_url}
									alt={team.name}
									class="w-16 h-16 object-contain"
								/>
							{/if}
							<div>
								<h2 class="text-xl font-bold">{team.name}</h2>
								<p class="text-theme-secondary">{team.abbreviation}</p>
							</div>
						</div>
					</a>
				{/each}
			</div>
			
			<p class="mt-8 text-theme-secondary text-sm">
				Total Teams: {teams.length}
			</p>
		{/if}
	</div>
</div>

