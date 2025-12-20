<script lang="ts">
	import type { BracketData } from '$lib/api';
	import BracketMatch from './BracketMatch.svelte';
	
	let { bracket, tournamentName }: { bracket: BracketData; tournamentName: string } = $props();
</script>

<div class="bracket-container overflow-x-auto pb-8">
	<div class="min-w-[1200px] p-8">
		<!-- Tournament Title -->
		<h2 class="text-2xl font-bold mb-8 text-center">{tournamentName}</h2>
		
		<div class="flex flex-col gap-16">
			<!-- Winners Bracket -->
			<div>
				<h3 class="text-lg font-semibold mb-6 text-theme-secondary uppercase tracking-wider">Winners Bracket</h3>
				<div class="flex gap-8 items-start">
					<!-- Winners Round 1 -->
					<div class="flex flex-col gap-4">
						<p class="text-xs text-theme-secondary uppercase tracking-wider mb-2">Winners Round 1</p>
						{#each bracket.winners_r1 as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.winners_r1.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
					
					<!-- Connector Lines -->
					<div class="flex items-center self-center">
						<div class="w-8 h-px bg-theme-secondary/50"></div>
					</div>
					
					<!-- Winners Round 2 -->
					<div class="flex flex-col gap-4 justify-center" style="margin-top: 2rem;">
						<p class="text-xs text-theme-secondary uppercase tracking-wider mb-2">Winners Round 2</p>
						{#each bracket.winners_r2 as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.winners_r2.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
					
					<!-- Connector Lines -->
					<div class="flex items-center self-center">
						<div class="w-8 h-px bg-theme-secondary/50"></div>
					</div>
					
					<!-- Winners Finals -->
					<div class="flex flex-col gap-4 justify-center" style="margin-top: 4rem;">
						<p class="text-xs text-theme-secondary uppercase tracking-wider mb-2">Winners Finals</p>
						{#each bracket.winners_finals as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.winners_finals.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
					
					<!-- Connector to Grand Finals -->
					<div class="flex items-center self-center">
						<div class="w-16 h-px bg-theme-secondary/50"></div>
					</div>
					
					<!-- Grand Finals -->
					<div class="flex flex-col gap-4 justify-center" style="margin-top: 4rem;">
						<p class="text-xs text-green-500 uppercase tracking-wider mb-2 font-semibold">Grand Finals</p>
						{#each bracket.grand_finals as match (match.id)}
							<div class="ring-2 ring-green-500/30 rounded-sm">
								<BracketMatch {match} />
							</div>
						{/each}
						{#if bracket.grand_finals.length === 0}
							<div class="w-52 h-16 bg-green-500/10 rounded-sm flex items-center justify-center text-green-500 text-sm ring-2 ring-green-500/30">
								TBD
							</div>
						{/if}
					</div>
				</div>
			</div>
			
			<!-- Elimination Bracket -->
			<div>
				<h3 class="text-lg font-semibold mb-6 text-theme-secondary uppercase tracking-wider">Elimination Bracket</h3>
				<div class="flex gap-8 items-start">
					<!-- Elimination Round 1 -->
					<div class="flex flex-col gap-4">
						<p class="text-xs text-theme-secondary uppercase tracking-wider mb-2">Elimination Round 1</p>
						{#each bracket.elim_r1 as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.elim_r1.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
					
					<!-- Connector -->
					<div class="flex items-center self-center">
						<div class="w-8 h-px bg-theme-secondary/50"></div>
					</div>
					
					<!-- Elimination Round 2 -->
					<div class="flex flex-col gap-4 justify-center" style="margin-top: 2rem;">
						<p class="text-xs text-theme-secondary uppercase tracking-wider mb-2">Elimination Round 2</p>
						{#each bracket.elim_r2 as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.elim_r2.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
					
					<!-- Connector -->
					<div class="flex items-center self-center">
						<div class="w-8 h-px bg-theme-secondary/50"></div>
					</div>
					
					<!-- Elimination Round 3 -->
					<div class="flex flex-col gap-4 justify-center" style="margin-top: 4rem;">
						<p class="text-xs text-theme-secondary uppercase tracking-wider mb-2">Elimination Round 3</p>
						{#each bracket.elim_r3 as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.elim_r3.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
					
					<!-- Connector -->
					<div class="flex items-center self-center">
						<div class="w-8 h-px bg-theme-secondary/50"></div>
					</div>
					
					<!-- Elimination Finals -->
					<div class="flex flex-col gap-4 justify-center" style="margin-top: 6rem;">
						<p class="text-xs text-red-500 uppercase tracking-wider mb-2">Elimination Finals</p>
						{#each bracket.elim_finals as match (match.id)}
							<BracketMatch {match} />
						{/each}
						{#if bracket.elim_finals.length === 0}
							<div class="w-52 h-16 bg-theme-secondary/30 rounded-sm flex items-center justify-center text-theme-secondary text-sm">
								TBD
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.bracket-container {
		background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-primary) 100%);
		border-radius: 0.5rem;
	}
</style>

