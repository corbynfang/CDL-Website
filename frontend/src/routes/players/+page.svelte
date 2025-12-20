<script lang="ts">
	import type { Player } from '$lib/api';
	
	// In Svelte 5, we get data from load() via $props()
	// This is equivalent to useLoaderData() in React Router
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	// This ensures the values update if data changes
	const players = $derived(data.players as Player[]);
	const error = $derived(data.error as string | null);
</script>

<svelte:head>
	<title>Players - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<h1 class="text-4xl font-bold mb-8 pb-4">PLAYERS</h1>
		
		{#if error}
			<p class="text-red-500">{error}</p>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="border-b border-theme">
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Gamertag
							</th>
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Name
							</th>
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Country
							</th>
							<th class="text-left py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Role
							</th>
							<th class="text-center py-4 text-theme-secondary text-sm uppercase tracking-wider">
								Status
							</th>
						</tr>
					</thead>
					<tbody>
						{#each players as player (player.id)}
							<tr class="border-b border-theme hover:bg-theme-secondary/50 transition-colors">
								<td class="py-4">
									<a 
										href="/players/{player.id}" 
										class="font-semibold hover:opacity-70 transition-opacity"
									>
										{player.gamertag}
									</a>
								</td>
								<td class="py-4 text-theme-secondary">
									{player.first_name && player.last_name
										? `${player.first_name} ${player.last_name}`
										: '—'}
								</td>
								<td class="py-4 text-theme-secondary">
									{player.country || '—'}
								</td>
								<td class="py-4 text-theme-secondary">
									{player.role || '—'}
								</td>
								<td class="py-4 text-center">
									<span class={player.is_active ? '' : 'text-theme-secondary'}>
										{player.is_active ? 'Active' : 'Inactive'}
									</span>
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

