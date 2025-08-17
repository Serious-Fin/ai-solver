import { toast } from 'svelte-sonner'
import { browser } from '$app/environment'
import { PUBLIC_API_BASE_URL_SERVER } from '$env/static/public'

export function handleFrontendError(msgToUser: string, err: Error) {
	toast.error(msgToUser)

	if (browser) {
		fetch('/api/log-error', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				message: err.message
			})
		}).catch((logErr) => {
			console.error('Failed to log error to server:', logErr)
		})
	}
}

export function showWarning(msg: string) {
	toast.warning(msg)
}

export function showSuccess(msg: string) {
	toast.success(msg)
}

const arrayStartRegex = /\[\][a-zA-Z0-9_]*{/g
const arrayEndRegex = /}/g
export function printHumanReadable(input: string): string {
	input = input.replaceAll(arrayStartRegex, '[')
	input = input.replaceAll(arrayEndRegex, ']')
	return input
}

export function getDifficultyName(difficulty: number): string {
	switch (difficulty) {
		case 1:
			return 'easy'
		case 2:
			return 'medium'
		case 3:
			return 'hard'
		default:
			return 'legendary'
	}
}

/**
 * This helper function is needed because based on whether API is called from browser or server,
 * different URLS are used.
 * @returns api host name
 */
export function getApiName(): string {
	if (browser) {
		return '/api/proxy'
	} else {
		return PUBLIC_API_BASE_URL_SERVER
	}
}
