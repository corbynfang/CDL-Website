<script lang="ts">
	// Home page - no data fetching needed here
	// In SvelteKit, page-level data comes from +page.ts or +page.server.ts
</script>

<svelte:head>
	<title>CDLYTICS - Call of Duty League Statistics</title>
	<meta name="description" content="Call of Duty League Statistics & Analytics" />
</svelte:head>

<div class="min-h-screen">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
		<!-- Hero Section -->
		<div class="text-center mb-16">
			<h1 class="text-5xl font-bold mb-4">CDLYTICS</h1>
			<p class="text-xl text-theme-secondary">
				Call of Duty League Statistics & Analytics
			</p>
		</div>
		
		<!-- Navigation Cards -->
		<div class="grid md:grid-cols-3 gap-8 max-w-4xl mx-auto">
			<a 
				href="/players" 
				class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200"
			>
				<h2 class="text-2xl font-bold mb-2">PLAYERS</h2>
				<p class="text-theme-secondary">
					View player statistics and performance data
				</p>
			</a>
			
			<a 
				href="/teams" 
				class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200"
			>
				<h2 class="text-2xl font-bold mb-2">TEAMS</h2>
				<p class="text-theme-secondary">
					Browse teams and roster information
				</p>
			</a>
			
			<a 
				href="/stats" 
				class="p-8 bg-theme-secondary shadow-md shadow-black/10 hover:shadow-lg hover:-translate-y-1 transition-all duration-200"
			>
				<h2 class="text-2xl font-bold mb-2">STATS</h2>
				<p class="text-theme-secondary">
					K/D leaderboards and rankings
				</p>
			</a>
		</div>
	</div>
</div>
