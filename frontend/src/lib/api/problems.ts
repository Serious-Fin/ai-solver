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

const BASE_URL = 'http://127.0.0.1:8080'

export async function getProblems(userId: number): Promise<Problem[]> {
	try {
		const response = await fetch(`${BASE_URL}/problems?user=${userId}`)
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
		const response = await fetch(`${BASE_URL}/problems/${problemId}?user=${userId}`)
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching problem with id '${problemId}' (userId: ${userId}) ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const problem: Problem = await response.json()

		const codeTemplate = await fetch(`${BASE_URL}/problems/${problemId}/go`)
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
