<script lang="ts">
	import type { BracketMatch } from '$lib/api';
	
	let { match }: { match: BracketMatch } = $props();
	
	// Determine winner styling
	const team1Won = $derived(match.winner_id === match.team1_id);
	const team2Won = $derived(match.winner_id === match.team2_id);
	const hasWinner = $derived(match.winner_id !== null);
</script>

<div class="bracket-match bg-theme-secondary rounded-sm overflow-hidden shadow-md w-52">
	<!-- Team 1 -->
	<div class="flex items-center justify-between px-3 py-2 border-b border-theme {team1Won ? 'bg-green-500/10' : hasWinner && !team1Won ? 'opacity-50' : ''}">
		<div class="flex items-center gap-2 flex-1 min-w-0">
			{#if match.team1_logo}
				<img 
					src={match.team1_logo} 
					alt={match.team1_abbr}
					class="w-5 h-5 object-contain flex-shrink-0"
				/>
			{/if}
			<span class="text-sm font-medium truncate {team1Won ? 'font-bold' : ''}">
				{match.team1_abbr || match.team1_name}
			</span>
		</div>
		<span class="text-sm font-bold ml-2 {team1Won ? 'text-green-500' : ''}">
			{match.team1_score}
			{#if team1Won}
				<span class="text-green-500 ml-1">◀</span>
			{/if}
		</span>
	</div>
	
	<!-- Team 2 -->
	<div class="flex items-center justify-between px-3 py-2 {team2Won ? 'bg-green-500/10' : hasWinner && !team2Won ? 'opacity-50' : ''}">
		<div class="flex items-center gap-2 flex-1 min-w-0">
			{#if match.team2_logo}
				<img 
					src={match.team2_logo} 
					alt={match.team2_abbr}
					class="w-5 h-5 object-contain flex-shrink-0"
				/>
			{/if}
			<span class="text-sm font-medium truncate {team2Won ? 'font-bold' : ''}">
				{match.team2_abbr || match.team2_name}
			</span>
		</div>
		<span class="text-sm font-bold ml-2 {team2Won ? 'text-green-500' : ''}">
			{match.team2_score}
			{#if team2Won}
				<span class="text-green-500 ml-1">◀</span>
			{/if}
		</span>
	</div>
</div>

