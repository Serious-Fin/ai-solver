import type { PageServerLoad } from './$types'
import { getProblems } from '$lib/api/problems'

export const load: PageServerLoad = async ({ locals }) => {
	return {
		problems: await getProblems(locals.user?.id ?? -1),
		user: locals.user
	}
}
