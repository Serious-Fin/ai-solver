import { toast } from 'svelte-sonner'
import { browser } from '$app/environment'
import { PUBLIC_API_BASE_URL as API_BASE_URL } from '$env/static/public'

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
 * Calling from server side - use "api" abbreviation because docker compose handles it
 * Calling from browser - use "localhost"
 * @returns api host name
 */
export function getApiName(): string {
	if (browser) {
		return 'http://localhost:8080'
	} else {
		return API_BASE_URL
	}
}
