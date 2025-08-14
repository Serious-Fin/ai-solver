import type { RequestHandler } from './$types'
import { OAuth2Client } from 'google-auth-library'
import { PUBLIC_GOOGLE_OAUTH_CLIENT_ID } from '$env/static/public'
import { GOOGLE_OAUTH_CLIENT_SECRET } from '$env/static/private'
import { redirect } from '@sveltejs/kit'
import { createSessionForUser, type User } from '$lib/api/users'

const client = new OAuth2Client({
	client_id: PUBLIC_GOOGLE_OAUTH_CLIENT_ID,
	client_secret: GOOGLE_OAUTH_CLIENT_SECRET
})

async function getUserInfo(token: string): Promise<User> {
	const ticket = await client.verifyIdToken({
		idToken: token,
		audience: PUBLIC_GOOGLE_OAUTH_CLIENT_ID
	})
	const payload = ticket.getPayload()
	if (!payload) {
		throw new Error('oauth payload is empty or has no email')
	}
	return {
		id: payload.sub,
		name: payload.name,
		profilePic: payload.picture,
		email: payload.email
	}
}

export const POST: RequestHandler = async ({ url, request, cookies }) => {
	const formData = await request.formData()
	const credential = formData.get('credential')
	if (!credential) {
		return new Response('No Google credential in OAuth response', { status: 401 })
	}
	let user: User
	try {
		user = await getUserInfo(credential as string)
	} catch (err) {
		return new Response('Could not get email from OAuth response', { status: 401 })
	}

	let sessionId: string
	try {
		sessionId = await createSessionForUser(user)
	} catch (err) {
		return new Response('Could not create session for user', { status: 401 })
	}

	cookies.set('session', sessionId, {
		httpOnly: true,
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7
	})

	const redirectTo = url.searchParams.get('redirectTo')
	const decodedRedirectTo = redirectTo ? decodeURIComponent(redirectTo) : '/'
	throw redirect(302, `${decodedRedirectTo}`)
}

// TODO: handle the edge cases appropriately in GET and POST methods
