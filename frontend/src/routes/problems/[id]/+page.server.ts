import type { PageServerLoad, Actions } from './$types'
import { getProblemById } from '$lib/api/problems'
import { query, type QueryRequest } from '$lib/api/query'
import { fail } from '@sveltejs/kit'

export const load: PageServerLoad = async ({ params }) => {
	return {
		problem: await getProblemById(params.id)
	}
}

export const actions = {
	query: async ({ request }) => {
		const data = await request.formData()
		const params: QueryRequest = {
			input: data.get('input') as string,
			code: data.get('code') as string,
			language: data.get('language') as string,
			agent: data.get('agent') as string,
			sessionId: data.get('sessionId') as string
		}
		try {
			const response = await query(params)
			return { response }
		} catch (err) {
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			})
		}
	}
} satisfies Actions
