import type { RequestHandler } from './$types'
import { PUBLIC_GITHUB_OAUTH_CLIENT_ID } from '$env/static/public'
import { GITHUB_OAUTH_CLIENT_SECRET } from '$env/static/private'

export const GET: RequestHandler = async ({ url, request, cookies }) => {
	if (url.searchParams.has('code')) {
		// handle received info
	} else {
		// make request to github
	}
	return new Response('No Google credential in OAuth response', { status: 401 })
}
