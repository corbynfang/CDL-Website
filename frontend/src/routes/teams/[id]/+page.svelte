<script lang="ts">
	import type { Team, Player } from '$lib/api';
	
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	const team = $derived(data.team as Team | null);
	const players = $derived(data.players as Player[]);
	const error = $derived(data.error as string | null);
</script>

<svelte:head>
	<title>{team?.name || 'Team'} - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<!-- Back Button -->
		<a 
			href="/teams" 
			class="text-theme-secondary hover:text-theme-primary mb-8 inline-block transition-colors"
		>
			← Back to Teams
		</a>
		
		{#if error || !team}
			<p class="text-red-500">{error || 'Team not found'}</p>
		{:else}
			<!-- Team Header -->
			<div class="pb-8 mb-8">
				<div class="flex items-center space-x-6">
					{#if team.logo_url}
						<img
							src={team.logo_url}
							alt={team.name}
							class="w-32 h-32 object-contain"
						/>
					{/if}
					<div>
						<h1 class="text-5xl font-bold mb-2">{team.name}</h1>
						<p class="text-xl text-theme-secondary">{team.abbreviation}</p>
					</div>
				</div>
			</div>
			
			<!-- Roster -->
			<div>
				<h2 class="text-2xl font-bold mb-6">ROSTER</h2>
				{#if players.length > 0}
					<div class="space-y-4">
						{#each players as player (player.id)}
							<a
								href="/players/{player.id}"
								class="p-6 bg-theme-secondary shadow-md shadow-black/10 flex items-center justify-between hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200"
							>
								<div class="flex items-center space-x-4">
									{#if player.avatar_url}
										<img
											src={player.avatar_url}
											alt={player.gamertag}
											class="w-12 h-12 object-cover"
										/>
									{/if}
									<div>
										<p class="font-bold text-lg">{player.gamertag}</p>
										<p class="text-sm text-theme-secondary">
											{player.first_name && player.last_name
												? `${player.first_name} ${player.last_name}`
												: '—'}
										</p>
									</div>
								</div>
								<div class="text-right">
									<p class="text-theme-secondary text-sm">{player.role || '—'}</p>
								</div>
							</a>
						{/each}
					</div>
				{:else}
					<p class="text-theme-secondary">No players found</p>
				{/if}
			</div>
		{/if}
	</div>
</div>

