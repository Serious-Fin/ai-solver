import type { Handle } from '@sveltejs/kit'

export const handle: Handle = async ({ event, resolve }) => {
	const sessionId = event.cookies.get('session')
	if (sessionId) {
		event.locals.sessionId = sessionId
	}
	return resolve(event)
}
