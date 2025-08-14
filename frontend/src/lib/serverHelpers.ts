import { DISCORD_TOKEN } from '$env/static/private'
import { DISCORD_CHANNEL_ID } from '$env/static/private'

// had to export this function here, because having it in the same
// file as some other method used in frontend, results in sveltekit error
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
