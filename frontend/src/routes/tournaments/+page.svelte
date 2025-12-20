<script lang="ts">
	import type { Tournament } from '$lib/api';
	
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	const tournaments = $derived(data.tournaments as Tournament[]);
	const error = $derived(data.error as string | null);
	
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
	<title>Tournaments - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<h1 class="text-4xl font-bold mb-8 pb-4">TOURNAMENTS</h1>
		
		{#if error}
			<p class="text-red-500">{error}</p>
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
			
			{#if tournaments.length === 0}
				<div class="bg-theme-secondary p-8 text-center">
					<p class="text-theme-secondary">No tournaments available</p>
				</div>
			{/if}
			
			<p class="mt-8 text-theme-secondary text-sm">
				Total Tournaments: {tournaments.length}
			</p>
		{/if}
	</div>
</div>

