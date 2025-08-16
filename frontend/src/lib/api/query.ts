import { getApiName } from '$lib/helpers'

export interface QueryRequest {
	input: string
	code: string
	language: string
	agent: string
	sessionId: string
}

export interface QueryResponse {
	response: string
}

export async function query(params: QueryRequest): Promise<string> {
	try {
		const response = await fetch(`${getApiName()}/query/${params.sessionId}`, {
			body: JSON.stringify({
				input: params.input,
				code: params.code,
				language: params.language,
				agent: params.agent
			}),
			headers: {
				'Content-Type': 'application/json'
			},
			method: 'POST'
		})
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error querying agent ${response.status} - ${errorBody.message || 'Unknown error'}`
			)
		}
		const parsedResp: QueryResponse = await response.json()
		return parsedResp.response
	} catch (err) {
		if (err instanceof Error) {
			throw Error(`Could not call query endpoint: ${JSON.stringify(err.message)}`)
		}
		throw err
	}
}
