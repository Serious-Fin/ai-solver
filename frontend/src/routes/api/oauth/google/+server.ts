import type { RequestHandler } from './$types'
import { OAuth2Client } from 'google-auth-library'
import { PUBLIC_GOOGLE_OAUTH_CLIENT_ID } from '$env/static/public'
import { GOOGLE_OAUTH_CLIENT_SECRET } from '$env/static/private'

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

	// TODO:
	// try and find a user with provided email in DB (return id or -1)
	// if none exists, create a user with provided email (return id)

	// try and get sessionId from session table by userId where it is not yet expired
	// if found, add one week ot expire (UPDATE) and set sessionId
	// if not found, delete all sessionIds where userID is set (clean up expired)
	// create new session id
	// set sessionId
	// redirect

	return new Response(null, { status: 200 })
}

/*
Needed operations:
GET /user/get?email=EMAIL
POST /user/create body: {email: EMAIL}

GET /session?userId=USERID
UPDATE /session body: {expireAt: NOW + 1 WEEK}
DELETE /session?userId=USERID
POST /session {userId: USERID}
*/
