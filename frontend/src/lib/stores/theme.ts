import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Check for saved preference or system preference
function getInitialTheme(): 'light' | 'dark' {
	if (!browser) return 'light';
	
	const saved = localStorage.getItem('theme');
	if (saved === 'dark' || saved === 'light') return saved;
	
	// Check system preference
	if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
		return 'dark';
	}
	
	return 'light';
}

// Create the theme store
function createThemeStore() {
	const { subscribe, set, update } = writable<'light' | 'dark'>(getInitialTheme());
	
	return {
		subscribe,
		toggle: () => {
			update(current => {
				const newTheme = current === 'light' ? 'dark' : 'light';
				if (browser) {
					localStorage.setItem('theme', newTheme);
					document.documentElement.classList.toggle('dark', newTheme === 'dark');
				}
				return newTheme;
			});
		},
		set: (value: 'light' | 'dark') => {
			if (browser) {
				localStorage.setItem('theme', value);
				document.documentElement.classList.toggle('dark', value === 'dark');
			}
			set(value);
		},
		init: () => {
			if (browser) {
				const theme = getInitialTheme();
				document.documentElement.classList.toggle('dark', theme === 'dark');
				set(theme);
			}
		}
	};
}

export const theme = createThemeStore();

