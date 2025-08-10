import devtoolsJson from 'vite-plugin-devtools-json'
import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

export default defineConfig({
	plugins: [sveltekit(), devtoolsJson()],
	server: {
		headers: {
			'Cross-Origin-Opener-Policy': 'same-origin-allow-popups'
		}
	}
})
