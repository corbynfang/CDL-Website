<script lang="ts">
	import '../app.css';
	import { theme } from '$lib/stores/theme';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
	import { onMount } from 'svelte';
	
	// Svelte 5 uses 'children' snippet for slot content
	let { children } = $props();
	
	// Initialize theme on mount (client-side only)
	onMount(() => {
		theme.init();
	});
</script>

<div class="min-h-screen bg-theme-primary text-theme-primary transition-colors duration-300">
	<!-- Header -->
	<header class="border-b border-theme">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex items-center justify-between h-16">
				<a 
					href="/" 
					class="text-xl font-bold tracking-tight hover:opacity-80 transition-opacity"
				>
					CDLYTICS
				</a>
				
				<nav class="flex items-center space-x-8">
					<a 
						href="/players" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Players
					</a>
					<a 
						href="/teams" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Teams
					</a>
					<a 
						href="/tournaments" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Brackets
					</a>
					<a 
						href="/stats" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Stats
					</a>
					<a 
						href="/transfers" 
						class="text-sm uppercase tracking-wider text-theme-secondary hover:text-theme-primary transition-colors"
					>
						Transfers
					</a>
					
					<!-- Theme Toggle Button -->
					<ThemeToggle />
				</nav>
			</div>
		</div>
	</header>
	
	<!-- Main Content - this is where page content goes (like React's <Outlet />) -->
	<main>
		{@render children()}
	</main>
	
	<!-- Footer -->
	<footer class="border-t border-theme mt-16">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<p class="text-center text-theme-secondary text-sm">CDLytics © 2025</p>
		</div>
	</footer>
</div>
