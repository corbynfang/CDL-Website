<script lang="ts">
	import type { PlayerStats } from '$lib/api';
	
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	const statsData = $derived(data.statsData);
	const error = $derived(data.error as string | null);
	
	// Get players array and sort by K/D
	const players = $derived(
		(statsData?.players || [])
			.slice()
			.sort((a: PlayerStats, b: PlayerStats) => b.season_kd - a.season_kd)
	);
</script>

<svelte:head>
	<title>K/D Leaderboard - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<h1 class="text-4xl font-bold mb-8 pb-4">K/D LEADERBOARD</h1>
		
		{#if error}
			<p class="text-red-500">{error}</p>
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

