import { toast } from 'svelte-sonner'
import { browser } from '$app/environment'
import { DISCORD_TOKEN } from '$env/static/private'
import { DISCORD_CHANNEL_ID } from '$env/static/private'

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

export async function sendToDiscord(message: string) {
	const chunkSize = 2000
	for (let i = 0; i < message.length; i += chunkSize) {
		const chunk = message.slice(i, i + chunkSize)
		const resp = await fetch(`https://discord.com/api/channels/${DISCORD_CHANNEL_ID}/messages`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: DISCORD_TOKEN
			},
			body: JSON.stringify({
				content: `frontend-msg: ${chunk}`
			})
		})

		if (!resp.ok) {
			console.error('Discord API Error:', resp.status, await resp.text())
		}
	}
}
