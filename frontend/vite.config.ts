import devtoolsJson from 'vite-plugin-devtools-json'
import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'
import { readFileSync } from 'node:fs'

const isProduction = process.env.NODE_ENV === 'production'

let httpsConfig = undefined
if (isProduction) {
	try {
		httpsConfig = {
			key: readFileSync('/etc/letsencrypt/live/data-gen.eu/privkey.pem'),
			cert: readFileSync('/etc/letsencrypt/live/data-gen.eu/fullchain.pem')
		}
	} catch (error) {
		console.warn('SSL certificates not found, running without HTTPS')
	}
}

export default defineConfig({
	plugins: [sveltekit(), devtoolsJson()],
	server: {
		headers: {
			'Cross-Origin-Opener-Policy': 'same-origin-allow-popups'
		},
		...(httpsConfig && { https: httpsConfig }),
		host: isProduction ? '0.0.0.0' : 'localhost',
		port: isProduction ? 443 : 5173
	}
})
