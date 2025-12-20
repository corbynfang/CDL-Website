<script lang="ts">
	import type { PlayerTransfer } from '$lib/api';
	
	let { data } = $props();
	
	// Use $derived for reactive access to data properties
	const transfers = $derived(data.transfers as PlayerTransfer[]);
	const error = $derived(data.error as string | null);
	
	// Format date for display
	function formatDate(dateString: string): string {
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
	
	// Get transfer type badge styles
	function getTransferTypeStyles(type: string): string {
		switch (type.toLowerCase()) {
			case 'signing':
				return 'bg-green-500/20 text-green-600 dark:text-green-400';
			case 'release':
				return 'bg-red-500/20 text-red-600 dark:text-red-400';
			case 'trade':
				return 'bg-blue-500/20 text-blue-600 dark:text-blue-400';
			default:
				return 'bg-gray-500/20 text-gray-600 dark:text-gray-400';
		}
	}
</script>

<svelte:head>
	<title>Transfers - CDLYTICS</title>
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<h1 class="text-4xl font-bold mb-8 pb-4">TRANSFERS</h1>
		
		{#if error}
			<p class="text-red-500">{error}</p>
		{:else}
			<div class="space-y-4">
				{#if transfers.length > 0}
					{#each transfers as transfer (transfer.id)}
						<div class="bg-theme-secondary p-6 shadow-md shadow-black/10">
							<div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
								<!-- Left: Player Info -->
								<div class="flex items-center space-x-4 flex-1">
									{#if transfer.player?.avatar_url}
										<img
											src={transfer.player.avatar_url}
											alt={transfer.player.gamertag}
											class="w-16 h-16 object-cover shadow-sm"
										/>
									{/if}
									<div class="flex-1">
										<div class="flex items-center gap-3 mb-2">
											<a
												href="/players/{transfer.player_id}"
												class="text-xl font-bold hover:opacity-70 transition-opacity"
											>
												{transfer.player?.gamertag || 'Unknown Player'}
											</a>
											<span class="px-3 py-1 text-xs font-medium {getTransferTypeStyles(transfer.transfer_type)}">
												{transfer.transfer_type.toUpperCase()}
											</span>
										</div>
										
										<!-- Transfer Path -->
										<div class="flex items-center gap-2 text-sm">
											<span class="text-theme-secondary">
												{transfer.from_team?.name || 'Free Agent'}
											</span>
											<span>→</span>
											<span class="font-semibold">
												{transfer.to_team?.name || 'Free Agent'}
											</span>
										</div>
									</div>
								</div>
								
								<!-- Right: Details -->
								<div class="flex flex-col md:items-end gap-2 text-sm">
									<div class="flex items-center gap-4 text-theme-secondary">
										{#if transfer.role}
											<span class="uppercase tracking-wider">{transfer.role}</span>
										{/if}
										{#if transfer.season}
											<span class="uppercase tracking-wider">{transfer.season}</span>
										{/if}
									</div>
									<p class="text-theme-secondary">
										{formatDate(transfer.transfer_date)}
									</p>
								</div>
							</div>
							
							{#if transfer.description}
								<div class="mt-4 pt-4 border-t border-theme">
									<p class="text-sm text-theme-secondary">{transfer.description}</p>
								</div>
							{/if}
						</div>
					{/each}
				{:else}
					<div class="bg-theme-secondary p-8 text-center">
						<p class="text-theme-secondary">No transfers available</p>
					</div>
				{/if}
			</div>
			
			<p class="mt-8 text-theme-secondary text-sm">
				Total Transfers: {transfers.length}
			</p>
		{/if}
	</div>
</div>

