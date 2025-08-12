import { getSession } from '$lib/api/users'
import type { Handle } from '@sveltejs/kit'

export const handle: Handle = async ({ event, resolve }) => {
	const sessionId = event.cookies.get('session')
	if (sessionId) {
		const session = await getSession(sessionId)
		event.locals.user = session.user
	}
	return resolve(event)
}
