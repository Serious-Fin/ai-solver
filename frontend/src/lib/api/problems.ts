import { getApiName } from '$lib/helpers'

export interface Problem {
	id: string
	title: string
	difficulty: number
	description?: string
	testCases?: TestCase[]
	goPlaceholder?: string
	isCompleted: boolean
}

export interface TestCase {
	id: number
	inputs: string[]
	output: string
}

export async function getProblems(userId: string): Promise<Problem[]> {
	try {
		const response = await fetch(`${getApiName()}/problems?user=${userId}`)
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problems ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const problems: Problem[] = await response.json()
		return problems
	} catch (error) {
		throw error
	}
}

export async function getProblemById(problemId: string, userId: string): Promise<Problem> {
	try {
		const response = await fetch(`${getApiName()}/problems/${problemId}?user=${userId}`)
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problem with id '${problemId}' (userId: ${userId}) ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const problem: Problem = await response.json()

		const codeTemplate = await fetch(`${getApiName()}/problems/${problemId}/go`)
		if (!codeTemplate.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problem template with id '${problemId}' ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		problem.goPlaceholder = await codeTemplate.json()
		return problem
	} catch (error) {
		throw error
	}
}

export async function markProblemCompleted(problemId: string, userId: string): Promise<void> {
	try {
		const response = await fetch(`${getApiName()}/problems/${problemId}`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				userId
			})
		})
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error marking problem as completed '${problemId}' (userId: ${userId}) ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
	} catch (err) {
		throw new Error(`Could not mark problem ${problemId} as completed: ${err}`)
	}
}
