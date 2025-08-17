import type { PageServerLoad, Actions } from './$types'
import { getProblemById, markProblemCompleted } from '$lib/api/problems'
import { query, type QueryRequest } from '$lib/api/query'
import { validate, type ValidateRequest } from '$lib/api/validate'
import { fail } from '@sveltejs/kit'

export const load: PageServerLoad = async ({ params, locals }) => {
	return {
		problem: await getProblemById(parseInt(params.id), locals.user?.id ?? '-1'),
		user: locals.user
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
	},
	runTests: async ({ request }) => {
		const data = await request.formData()
		const params: ValidateRequest = {
			code: data.get('code') as string,
			problemId: parseInt(data.get('problemId') as string),
			language: data.get('language') as string
		}
		try {
			const response = await validate(params)
			return { response }
		} catch (err) {
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			})
		}
	},
	markProblemCompleted: async ({ request }) => {
		const data = await request.formData()
		const problemId: number = parseInt(data.get('problemId') as string)
		const userId: string = data.get('userId') as string
		try {
			const response = await markProblemCompleted(problemId, userId)
			return { response }
		} catch (err) {
			return fail(500, {
				message: err instanceof Error ? err.message : 'Unknown error'
			})
		}
	}
} satisfies Actions
