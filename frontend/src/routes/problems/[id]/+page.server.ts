import type { PageServerLoad, Actions } from './$types'
import { getProblemById } from '$lib/api/problems'
import { query, type QueryRequest } from '$lib/api/query'

export const load: PageServerLoad = async ({ params }) => {
    return {
        problem: await getProblemById(params.id)
    }
}

export const actions = {
    query: async ({ request }) => {
        const data = await request.formData()
        const params: QueryRequest = {
            input: data.get("input") as string,
            code: data.get("code") as string,
            language: data.get("language") as string,
            agent: data.get("agent") as string,
            sessionId: data.get("sessionId") as string
        }
        return { success: true, response: await query(params) }
    }
} satisfies Actions
