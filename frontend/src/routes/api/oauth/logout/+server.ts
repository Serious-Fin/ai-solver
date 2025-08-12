import type { RequestHandler } from './$types'

export const POST: RequestHandler = async ({ cookies }) => {
	cookies.set('session', '', {
		httpOnly: true,
		sameSite: 'lax',
		path: '/',
		maxAge: 0
	})

	return new Response(null, { status: 204 })
}
