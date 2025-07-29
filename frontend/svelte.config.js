import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter(),
		csp: {
			mode: 'nonce', // Use nonces for both scripts and styles
			directives: {
				'default-src': ['self'],
				'script-src': ['self'], // SvelteKit will automatically add nonces
				'style-src': ['self', 'unsafe-inline'], // unsafe-inline needed for Svelte transitions
				'img-src': ['self', 'data:', 'https:', 'blob:'], // Added blob: for EPUB covers
				'font-src': ['self', 'data:', 'blob:'], // Added blob: for EPUB fonts
				'connect-src': ['self', 'blob:'], // Added blob: for EPUB resources
				'media-src': ['self', 'blob:'], // Added blob: for EPUB media
				'object-src': ['none'],
				'base-uri': ['self'],
				'frame-ancestors': ['none'],
				'form-action': ['self'],
				'upgrade-insecure-requests': true
			},
		}
	}
};

export default config;