const BASE_URL = 'http://127.0.0.1:8080'

export async function login(email: string): Promise<number> {
	try {
		const response = await fetch(`${BASE_URL}/login`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				email
			})
		})
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(`Error logging in ${response.status} - ${errorBody || 'Unknown error'}`)
		}
		const jsonResponse: { userId: number } = await response.json()
		return jsonResponse.userId
	} catch (error) {
		console.log('Error logging in', error)
		throw error
	}
}

export async function startSession(userId: number): Promise<string> {
	try {
		const response = await fetch(`${BASE_URL}/session`, {
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
			throw new Error(`Error starting session ${response.status} - ${errorBody || 'Unknown error'}`)
		}
		const jsonResponse: { sessionId: string } = await response.json()
		return jsonResponse.sessionId
	} catch (error) {
		console.log(`Error starting session`, error)
		throw error
	}
}
