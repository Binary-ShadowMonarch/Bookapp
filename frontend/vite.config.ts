import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { enhancedImages } from '@sveltejs/enhanced-img';

export default defineConfig({
	server: {
		host: '0.0.0.0',
		port: 3000,
		proxy: {
			'/api': {
				target: 'http://backend:8080',
				changeOrigin: true,
				ws: true
			}
		}
	},
	plugins: [tailwindcss(), enhancedImages(), sveltekit()]
});
