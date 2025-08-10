import type { RequestHandler } from './$types'
import { OAuth2Client } from 'google-auth-library'
import { PUBLIC_GOOGLE_OAUTH_CLIENT_ID } from '$env/static/public'
import { GOOGLE_OAUTH_CLIENT_SECRET } from '$env/static/private'
import { login, startSession } from '$lib/api/users'
import { redirect } from '@sveltejs/kit'

const client = new OAuth2Client({
	client_id: PUBLIC_GOOGLE_OAUTH_CLIENT_ID,
	client_secret: GOOGLE_OAUTH_CLIENT_SECRET
})

async function getEmail(token: string): Promise<string> {
	const ticket = await client.verifyIdToken({
		idToken: token,
		audience: PUBLIC_GOOGLE_OAUTH_CLIENT_ID
	})
	const payload = ticket.getPayload()
	if (!payload || !payload.email) {
		throw new Error('oauth payload is empty or has no email')
	}
	return payload.email
}

export const POST: RequestHandler = async ({ request, cookies }) => {
	const formData = await request.formData()
	const credential = formData.get('credential')
	if (!credential) {
		return new Response('No Google credential in OAuth response', { status: 401 })
	}
	let email: string
	try {
		email = await getEmail(credential as string)
	} catch (err) {
		return new Response('Could not get email from OAuth response', { status: 401 })
	}

	let sessionId: string
	try {
		const userId = await login(email)
		sessionId = await startSession(userId)
	} catch (err) {
		return new Response(`Could not sign user in: ${JSON.stringify(err)}`, { status: 401 })
	}

	cookies.set('session', sessionId, {
		httpOnly: true,
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7 // 1 week
	})

	throw redirect(302, '/problems')
}
