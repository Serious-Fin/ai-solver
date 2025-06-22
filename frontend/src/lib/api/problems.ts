export interface Problem {
    id: string;
    title: string;
    description: string;
    testCases: TestCase[];
}

export interface TestCase {
    input: string;
    output: string;
}

const BASE_URL = "http://127.0.0.1:8080"

export async function getProblems(): Promise<Problem[]> {
    try {
        const response = await fetch(`${BASE_URL}/problems`)
        if (!response.ok) {
            const errorBody = await response.json().catch(() => ({ message: response.statusText }))
            throw new Error(`Error fetching problems ${response.status} - ${errorBody || "Unknown error"}`)
        }
        const problems: Problem[] = await response.json()
        return problems
    } catch (error) {
        console.log("Error fetching problems", error)
        throw error
    }
}