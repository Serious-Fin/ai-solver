import type { PageServerLoad } from './$types'
import { getProblems } from '$lib/api/problems'

export const load: PageServerLoad = async ({ locals }) => {
	return {
		problems: await getProblems(),
		sessionId: locals.sessionId
	}
}
