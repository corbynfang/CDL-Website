<script lang="ts">
	import type { Tournament, BracketResponse } from '$lib/api';
	import Bracket from '$lib/components/Bracket.svelte';
	
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	const tournament = $derived(data.tournament as Tournament | null);
	const bracketData = $derived(data.bracketData as BracketResponse | null);
	const error = $derived(data.error as string | null);
	
	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		try {
			return new Date(dateString).toLocaleDateString('en-US', {
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
	<title>{tournament?.name || 'Tournament'} - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<!-- Back Button -->
		<a 
			href="/tournaments" 
			class="text-theme-secondary hover:text-theme-primary mb-8 inline-block transition-colors"
		>
			← Back to Tournaments
		</a>
		
		{#if error || !tournament}
			<p class="text-red-500">{error || 'Tournament not found'}</p>
		{:else}
		
			<!-- Tournament Header -->
			<div class="mb-8">
				<h1 class="text-4xl font-bold mb-2">{tournament.name}</h1>
				<div class="flex flex-wrap gap-4 text-theme-secondary">
					{#if tournament.tournament_type}
						<span class="bg-theme-secondary px-3 py-1 rounded-full text-sm">
							{tournament.tournament_type}
						</span>
					{/if}
					{#if tournament.location}
						<span>📍 {tournament.location}</span>
					{/if}
					<span>📅 {formatDate(tournament.start_date)}</span>
					{#if tournament.prize_pool}
						<span>💰 ${tournament.prize_pool.toLocaleString()}</span>
					{/if}
				</div>
			</div>
			
			<!-- Bracket -->
			{#if bracketData?.bracket}
				<Bracket 
					bracket={bracketData.bracket} 
					tournamentName={bracketData.tournament_name} 
				/>
				
				<p class="mt-4 text-theme-secondary text-sm">
					Total Matches: {bracketData.total_matches}
				</p>
			{:else}
				<div class="bg-theme-secondary p-8 text-center rounded-lg">
					<p class="text-theme-secondary">No bracket data available for this tournament.</p>
					<p class="text-sm text-theme-secondary mt-2">
						Bracket data needs to be populated with match information including bracket_round and bracket_position fields.
					</p>
				</div>
			{/if}
		{/if}
	</div>
</div>

