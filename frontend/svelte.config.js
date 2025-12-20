import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),

	kit: {
		// Static adapter for deployment with Go backend
		adapter: adapter({
			pages: 'dist',      // Output to dist/ (same as React Vite)
			assets: 'dist',
			fallback: 'index.html',  // SPA fallback for dynamic routes
			precompress: false,
			strict: false  // Allow dynamic routes to be handled by fallback
		}),
		prerender: {
			handleHttpError: 'warn',
			handleUnseenRoutes: 'ignore'  // Dynamic routes handled by SPA fallback
		}
	}
};

export default config;
