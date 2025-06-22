import type { PageServerLoad } from './$types'
import { getProblemById } from '$lib/api/problems'

export const load: PageServerLoad = async ({ params }) => {
    return {
        problem: await getProblemById(params.id)
    }
}