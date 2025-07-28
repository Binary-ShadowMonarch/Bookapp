import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter(),
		csp: {
			directives: {
				'script-src': ['self'],
				'img-src': ['self', 'data:', 'https:', 'blob:'], // Added blob: support
			},
		}
	}
};

export default config;