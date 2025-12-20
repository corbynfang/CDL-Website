<script lang="ts">
	import { theme } from '$lib/stores/theme';
	
	// In Svelte 5, we use $props() for props and $state() for local state
	// The $theme syntax auto-subscribes to stores
</script>

<button
	onclick={() => theme.toggle()}
	class="p-2 rounded-lg bg-theme-secondary hover:opacity-80 transition-all duration-200"
	aria-label="Toggle theme"
>
	{#if $theme === 'light'}
		<!-- Moon icon for dark mode -->
		<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
			<path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
		</svg>
	{:else}
		<!-- Sun icon for light mode -->
		<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
			<path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
		</svg>
	{/if}
</button>

