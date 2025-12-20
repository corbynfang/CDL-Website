<script lang="ts">
	import type { Team } from '$lib/api';
	
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	const teams = $derived(data.teams as Team[]);
	const error = $derived(data.error as string | null);
</script>

<svelte:head>
	<title>Teams - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<h1 class="text-4xl font-bold mb-8 pb-4">TEAMS</h1>
		
		{#if error}
			<p class="text-red-500">{error}</p>
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

