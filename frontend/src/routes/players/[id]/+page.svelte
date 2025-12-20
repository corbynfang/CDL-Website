<script lang="ts">
	import type { Player, PlayerKDResponse, PlayerMatchesResponse } from '$lib/api';
	
	// Get data from load function
	let { data } = $props();
	
	// Svelte 5 reactive state - similar to React useState
	let activeTab = $state('last5');
	
	// Use $derived for reactive access to data properties
	const player = $derived(data.player as Player | null);
	const stats = $derived(data.stats as PlayerKDResponse | null);
	const matchesData = $derived(data.matches as PlayerMatchesResponse | null);
	const error = $derived(data.error as string | null);
	
	// Computed values (Svelte 5 uses $derived)
	const events = $derived(matchesData?.events || []);
	const allMatches = $derived(events.flatMap((event: any) => event.matches || []));
	const last5Matches = $derived(allMatches.slice(0, 5));
	const tournamentStats = $derived(stats?.tournament_stats || []);
	
	const overallKD = $derived(stats?.avg_kd || 0);
	const hpKD = $derived(stats?.hp_kd_ratio || 0);
	const sndKD = $derived(stats?.snd_kd_ratio || 0);
	const ctlKD = $derived(stats?.control_kd_ratio || 0);
	
	// Helper function for progress bar
	function getKDProgress(kd: number): number {
		if (!kd || kd <= 0) return 0;
		return Math.min((kd / 2.0) * 100, 100);
	}
	
	// Format birthdate
	function formatBirthdate(dateString?: string): string {
		if (!dateString) return '—';
		try {
			const date = new Date(dateString);
			return date.toLocaleDateString('en-US', {
				year: 'numeric',
				month: 'long',
				day: 'numeric'
			});
		} catch {
			return dateString;
		}
	}
</script>

<svelte:head>
	<title>{player?.gamertag || 'Player'} - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<!-- Back Button -->
		<a 
			href="/players" 
			class="text-theme-secondary hover:text-theme-primary mb-8 inline-block transition-colors"
		>
			← Back to Players
		</a>
		
		{#if error || !player}
			<p class="text-red-500">{error || 'Player not found'}</p>
		{:else}
			<!-- Top Section: Player Image (Left) and KD Stats Box (Right) -->
			<div class="grid md:grid-cols-2 gap-8 mb-8">
				<!-- Left: Player Image & Info -->
				<div class="flex flex-col items-center md:items-start">
					{#if player.avatar_url}
						<img
							src={player.avatar_url}
							alt={player.gamertag}
							class="w-48 h-48 object-cover shadow-md shadow-black/10 mb-6"
						/>
					{/if}
					
					<!-- Player Info Card -->
					<div class="w-full bg-theme-secondary p-6 shadow-md shadow-black/10">
						<div class="space-y-4">
							<div>
								<p class="text-xs uppercase tracking-wider text-theme-secondary mb-1">Name</p>
								<p class="text-lg font-bold">
									{player.first_name && player.last_name
										? `${player.first_name} ${player.last_name}`
										: player.gamertag}
								</p>
							</div>
							{#if player.country}
								<div>
									<p class="text-xs uppercase tracking-wider text-theme-secondary mb-1">Country</p>
									<p class="text-base">{player.country}</p>
								</div>
							{/if}
							{#if player.birthdate}
								<div>
									<p class="text-xs uppercase tracking-wider text-theme-secondary mb-1">Birthday</p>
									<p class="text-base">{formatBirthdate(player.birthdate)}</p>
								</div>
							{/if}
							{#if player.role}
								<div>
									<p class="text-xs uppercase tracking-wider text-theme-secondary mb-1">Role</p>
									<p class="text-base">{player.role}</p>
								</div>
							{/if}
							<div>
								<p class="text-xs uppercase tracking-wider text-theme-secondary mb-1">Status</p>
								<p class="text-base">{player.is_active ? 'Active' : 'Inactive'}</p>
							</div>
						</div>
					</div>
				</div>
				
				<!-- Right: KD Stats Box with Progress Bars -->
				<div class="bg-theme-secondary p-8 shadow-md shadow-black/10">
					<h2 class="text-2xl font-bold mb-6">K/D STATISTICS</h2>
					
					<div class="space-y-6">
						<!-- Overall KD -->
						<div>
							<div class="flex justify-between items-center mb-2">
								<p class="text-xs text-theme-secondary font-medium">Overall K/D</p>
								<p class="text-lg font-bold">{overallKD.toFixed(2)}</p>
							</div>
							<div class="w-full bg-gray-300 dark:bg-gray-700 h-3 overflow-hidden">
								<div
									class="bg-black dark:bg-white h-full transition-all duration-500"
									style="width: {getKDProgress(overallKD)}%"
								></div>
							</div>
						</div>
						
						<!-- Hardpoint KD -->
						<div>
							<div class="flex justify-between items-center mb-2">
								<p class="text-xs text-theme-secondary font-medium">Hardpoint K/D</p>
								<p class="text-lg font-bold">{hpKD.toFixed(2)}</p>
							</div>
							<div class="w-full bg-gray-300 dark:bg-gray-700 h-3 overflow-hidden">
								<div
									class="bg-black dark:bg-white h-full transition-all duration-500"
									style="width: {getKDProgress(hpKD)}%"
								></div>
							</div>
						</div>
						
						<!-- S&D KD -->
						<div>
							<div class="flex justify-between items-center mb-2">
								<p class="text-xs text-theme-secondary font-medium">S&D K/D</p>
								<p class="text-lg font-bold">{sndKD.toFixed(2)}</p>
							</div>
							<div class="w-full bg-gray-300 dark:bg-gray-700 h-3 overflow-hidden">
								<div
									class="bg-black dark:bg-white h-full transition-all duration-500"
									style="width: {getKDProgress(sndKD)}%"
								></div>
							</div>
						</div>
						
						<!-- Control KD -->
						<div>
							<div class="flex justify-between items-center mb-2">
								<p class="text-xs text-theme-secondary font-medium">Control K/D</p>
								<p class="text-lg font-bold">{ctlKD.toFixed(2)}</p>
							</div>
							<div class="w-full bg-gray-300 dark:bg-gray-700 h-3 overflow-hidden">
								<div
									class="bg-black dark:bg-white h-full transition-all duration-500"
									style="width: {getKDProgress(ctlKD)}%"
								></div>
							</div>
						</div>
					</div>
					
					<!-- Additional Stats Summary -->
					<div class="mt-8 pt-8 border-t border-theme grid grid-cols-3 gap-4">
						<div class="text-center">
							<p class="text-xs text-theme-secondary uppercase mb-1">Kills</p>
							<p class="text-xl font-bold">{stats?.total_kills?.toLocaleString() || '0'}</p>
						</div>
						<div class="text-center">
							<p class="text-xs text-theme-secondary uppercase mb-1">Deaths</p>
							<p class="text-xl font-bold">{stats?.total_deaths?.toLocaleString() || '0'}</p>
						</div>
						<div class="text-center">
							<p class="text-xs text-theme-secondary uppercase mb-1">Assists</p>
							<p class="text-xl font-bold">{stats?.total_assists?.toLocaleString() || '0'}</p>
						</div>
					</div>
				</div>
			</div>
			
			<!-- Bottom Section: Tabs and Content -->
			<div class="mt-12">
				<!-- Tabs -->
				<div class="flex space-x-1 border-b border-theme mb-8">
					<button
						onclick={() => activeTab = 'last5'}
						class="px-6 py-3 text-sm font-medium transition-all {activeTab === 'last5' 
							? 'border-b-2 border-current' 
							: 'text-theme-secondary'}"
					>
						Last 5 Matches
					</button>
					<button
						onclick={() => activeTab = 'matches'}
						class="px-6 py-3 text-sm font-medium transition-all {activeTab === 'matches' 
							? 'border-b-2 border-current' 
							: 'text-theme-secondary'}"
					>
						Matches
					</button>
					<button
						onclick={() => activeTab = 'eventStats'}
						class="px-6 py-3 text-sm font-medium transition-all {activeTab === 'eventStats' 
							? 'border-b-2 border-current' 
							: 'text-theme-secondary'}"
					>
						Event Stats
					</button>
				</div>
				
				<!-- Tab Content -->
				<div class="bg-theme-secondary p-8 shadow-md shadow-black/10">
					{#if activeTab === 'last5'}
						<h3 class="text-xl font-bold mb-6">Last 5 Matches</h3>
						{#if last5Matches.length > 0}
							<div class="space-y-4">
								{#each last5Matches as match, index}
									<div class="bg-theme-primary p-6">
										<div class="flex justify-between items-center">
											<div class="flex items-center space-x-4">
												<span class="text-sm font-bold w-8 text-center {match.result?.startsWith('W') ? '' : 'text-theme-secondary'}">
													{match.result?.charAt(0) || '—'}
												</span>
												<div>
													<p class="font-semibold">vs {match.opponent_abbr || match.opponent || 'Unknown'}</p>
													<p class="text-sm text-theme-secondary">
														{match.date ? new Date(match.date).toLocaleDateString() : '—'}
													</p>
												</div>
											</div>
											<div class="flex space-x-6 text-right">
												<div>
													<p class="text-xs text-theme-secondary uppercase">K/D</p>
													<p class="text-lg font-bold">
														{typeof match.kd === 'number' ? match.kd.toFixed(2) : '0.00'}
													</p>
												</div>
												<div>
													<p class="text-xs text-theme-secondary uppercase">Kills</p>
													<p class="text-lg font-bold">{match.kills || '0'}</p>
												</div>
												<div>
													<p class="text-xs text-theme-secondary uppercase">Deaths</p>
													<p class="text-lg font-bold">{match.deaths || '0'}</p>
												</div>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-theme-secondary">No matches available</p>
						{/if}
					{:else if activeTab === 'matches'}
						{#if events.length > 0}
							<div class="space-y-8">
								{#each events as event, eventIndex}
									<div>
										<h3 class="text-xl font-bold mb-4">{event.event} {event.year}</h3>
										{#if event.matches && event.matches.length > 0}
											<div class="overflow-x-auto">
												<table class="w-full">
													<thead>
														<tr class="border-b border-theme">
															<th class="text-left py-3 text-theme-secondary text-xs uppercase tracking-wider">Date</th>
															<th class="text-left py-3 text-theme-secondary text-xs uppercase tracking-wider">Opponent</th>
															<th class="text-left py-3 text-theme-secondary text-xs uppercase tracking-wider">Result</th>
															<th class="text-right py-3 text-theme-secondary text-xs uppercase tracking-wider">KD</th>
															<th class="text-right py-3 text-theme-secondary text-xs uppercase tracking-wider">K</th>
															<th class="text-right py-3 text-theme-secondary text-xs uppercase tracking-wider">D</th>
														</tr>
													</thead>
													<tbody>
														{#each event.matches as match}
															<tr class="border-b border-theme">
																<td class="py-3 text-theme-secondary text-sm">
																	{match.date ? new Date(match.date).toLocaleDateString() : '—'}
																</td>
																<td class="py-3 text-sm font-medium">
																	{match.opponent_abbr || match.opponent || '—'}
																</td>
																<td class="py-3 text-sm font-medium">{match.result || '—'}</td>
																<td class="py-3 text-right text-sm font-bold">
																	{typeof match.kd === 'number' ? match.kd.toFixed(2) : '0.00'}
																</td>
																<td class="py-3 text-right text-theme-secondary text-sm">{match.kills || '0'}</td>
																<td class="py-3 text-right text-theme-secondary text-sm">{match.deaths || '0'}</td>
															</tr>
														{/each}
													</tbody>
												</table>
											</div>
										{:else}
											<p class="text-theme-secondary">No matches for this event</p>
										{/if}
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-theme-secondary">No match data available</p>
						{/if}
					{:else if activeTab === 'eventStats'}
						<h3 class="text-xl font-bold mb-6">Event Statistics</h3>
						{#if tournamentStats.length > 0}
							<div class="space-y-4">
								{#each tournamentStats as tournament}
									<div class="bg-theme-primary p-6">
										<div class="flex justify-between items-start mb-4">
											<div>
												<p class="font-bold text-lg">{tournament.tournament_name || 'Tournament'}</p>
												<p class="text-sm text-theme-secondary">
													{tournament.matches || 0} matches • {tournament.maps_played || 0} maps
												</p>
											</div>
											<div class="text-right">
												<p class="text-xs text-theme-secondary uppercase mb-1">K/D</p>
												<p class="text-2xl font-bold">{tournament.kd_ratio?.toFixed(2) || '0.00'}</p>
											</div>
										</div>
										<div class="grid grid-cols-3 gap-4 text-center pt-4 border-t border-theme">
											<div>
												<p class="text-xs text-theme-secondary uppercase mb-1">Kills</p>
												<p class="text-lg font-bold">{tournament.kills?.toLocaleString() || '0'}</p>
											</div>
											<div>
												<p class="text-xs text-theme-secondary uppercase mb-1">Deaths</p>
												<p class="text-lg font-bold">{tournament.deaths?.toLocaleString() || '0'}</p>
											</div>
											<div>
												<p class="text-xs text-theme-secondary uppercase mb-1">Assists</p>
												<p class="text-lg font-bold">{tournament.assists?.toLocaleString() || '0'}</p>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-theme-secondary">No event statistics available</p>
						{/if}
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>

