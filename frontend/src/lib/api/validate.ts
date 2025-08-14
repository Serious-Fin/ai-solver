import type { TestRunOutput } from '$lib/TestStatusReporter'
import { PUBLIC_API_BASE_URL as API_BASE_URL } from '$env/static/public'

export interface ValidateRequest {
	problemId: string
	code: string
	language: string
}

export async function validate(req: ValidateRequest): Promise<TestRunOutput> {
	try {
		const resp = await fetch(`${API_BASE_URL}/validate`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(req)
		})
		if (!resp.ok) {
			const errorBody = await resp.json().catch(() => ({ message: resp.statusText }))
			throw new Error(
				`Error running tests ${resp.status} - ${errorBody.message || 'Unknown error'}`
			)
		}
		const testRunOutput: TestRunOutput = await resp.json()
		return testRunOutput
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call validate endpoint: ${JSON.stringify(err.message)}`)
		}
		throw err
	}
}
