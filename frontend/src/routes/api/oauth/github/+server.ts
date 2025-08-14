import type { RequestHandler } from './$types'
import { PUBLIC_GITHUB_OAUTH_CLIENT_ID, PUBLIC_FRONTEND_BASE_URL } from '$env/static/public'
import { GITHUB_OAUTH_CLIENT_SECRET } from '$env/static/private'
import { v4 as uuidv4 } from 'uuid'
import { error, redirect } from '@sveltejs/kit'
import { handleError } from '$lib/helpers'
import { createSessionForUser } from '$lib/api/users'

const ghStateCookieName = 'gh_state'

export const GET: RequestHandler = async ({ url, request, cookies }) => {
	const redirectToAfterLogin = url.searchParams.get('redirectTo') ?? '/'
	const redirectUri = `${PUBLIC_FRONTEND_BASE_URL}/api/oauth/github?redirectTo=${redirectToAfterLogin}`

	if (url.searchParams.has('code')) {
		const storedState = cookies.get(ghStateCookieName)
		const returnedState = url.searchParams.get('state')
		if (!storedState || storedState !== returnedState) {
			redirect(308, '/login')
		}

		let sessionId: string
		try {
			const code = url.searchParams.get('code') ?? ''
			const accessToken = await getAccessToken(code, redirectUri)
			const email = await getUserEmail(accessToken)
			console.log(email)
			sessionId = await createSessionForUser(email)
		} catch (err) {
			return new Response(`Could not create session for user (auth via github): ${err}`, {
				status: 401
			})
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
	} else {
		const scopes = 'read:user user:email'
		const state = uuidv4()
		const githubOAuthUrl =
			`https://github.com/login/oauth/authorize` +
			`?client_id=${encodeURIComponent(PUBLIC_GITHUB_OAUTH_CLIENT_ID)}` +
			`&redirect_uri=${encodeURIComponent(redirectUri)}` +
			`&scopes=${encodeURIComponent(scopes)}` +
			`&state=${encodeURIComponent(state)}`

		cookies.set(ghStateCookieName, state, {
			httpOnly: true,
			sameSite: 'lax',
			path: '/',
			maxAge: 60 * 10
		})
		redirect(307, githubOAuthUrl)
	}
}

async function getAccessToken(code: string, redirectUri: string): Promise<string> {
	try {
		const response = await fetch('https://github.com/login/oauth/access_token', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json'
			},
			body: JSON.stringify({
				client_id: PUBLIC_GITHUB_OAUTH_CLIENT_ID,
				client_secret: GITHUB_OAUTH_CLIENT_SECRET,
				code,
				redirect_uri: redirectUri
			})
		})
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error receiving github access token ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const { access_token } = await response.json()
		return access_token
	} catch (err) {
		throw new Error(`github access_token request failed: ${err}`)
	}
}

interface GithubEmailRecord {
	email: string
	primary: boolean
}

async function getUserEmail(accessToken: string): Promise<string> {
	try {
		const response = await fetch('https://api.github.com/user/emails', {
			headers: {
				Authorization: `Bearer ${accessToken}`,
				Accept: 'application/vnd.github+json'
			}
		})
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`response was non 2xx ${response.status} - ${errorBody.message || 'Unknown error'}`
			)
		}
		const allEmails: GithubEmailRecord[] = await response.json()
		if (allEmails.length === 0) {
			throw new Error('Github user does not have an email set')
		}
		const primaryEmail = allEmails.find((e) => e.primary)?.email
		if (!primaryEmail) {
			return allEmails[0].email
		}
		return primaryEmail
	} catch (err) {
		throw new Error(`could not get github user email: ${err}`)
	}
}
