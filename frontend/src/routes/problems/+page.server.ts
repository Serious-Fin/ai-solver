import type { PageServerLoad } from './$types'
import { getProblems } from '$lib/api/problems'

export const load: PageServerLoad = async () => {
	return {
		problems: await getProblems()
	}
}
