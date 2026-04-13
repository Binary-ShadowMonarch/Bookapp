import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { enhancedImages } from '@sveltejs/enhanced-img';

export default defineConfig({
	server: {
		allowedHosts: ['90f1f71e6b7c.ngrok-free.app']
	},
	plugins: [tailwindcss(), enhancedImages(), sveltekit()]
});
