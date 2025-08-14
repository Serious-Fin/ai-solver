import { PUBLIC_API_BASE_URL as API_BASE_URL } from '$env/static/public'

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

export async function getProblems(userId: number): Promise<Problem[]> {
	try {
		const response = await fetch(`${API_BASE_URL}/problems?user=${userId}`)
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problems ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const problems: Problem[] = await response.json()
		return problems
	} catch (error) {
		console.log('Error fetching problems', error)
		throw error
	}
}

export async function getProblemById(problemId: string, userId: number): Promise<Problem> {
	try {
		const response = await fetch(`${API_BASE_URL}/problems/${problemId}?user=${userId}`)
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problem with id '${problemId}' (userId: ${userId}) ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const problem: Problem = await response.json()

		const codeTemplate = await fetch(`${API_BASE_URL}/problems/${problemId}/go`)
		if (!codeTemplate.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problem template with id '${problemId}' ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		problem.goPlaceholder = await codeTemplate.json()
		return problem
	} catch (error) {
		console.log(`Error fetching problem with id '${problemId}' (userId: ${userId})`, error)
		throw error
	}
}

// TODO: passing all tests should mark the task as complete if user is signed in
