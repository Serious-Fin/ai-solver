import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'
import { sendToDiscord } from '$lib/serverHelpers'

export const POST: RequestHandler = async ({ request }) => {
	try {
		const { message } = await request.json()
		await sendToDiscord(message)
		return json({ success: true })
	} catch (err) {
		console.error('Failed to log to Discord:', err)
		return json({ success: false }, { status: 500 })
	}
}
