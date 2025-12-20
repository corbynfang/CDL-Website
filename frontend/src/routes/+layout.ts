// Disable SSR globally - all pages render client-side only
// This is needed for SPA mode with static adapter
export const ssr = false;
// Don't prerender - use SPA fallback for all routes
export const prerender = false;

