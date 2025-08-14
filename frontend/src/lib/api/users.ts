import { PUBLIC_API_BASE_URL as API_BASE_URL } from '$env/static/public'

export interface User {
	id: number
	email: string
}

export interface SessionInfo {
	user?: User
}

export async function createSessionForUser(email: string): Promise<string> {
	let sessionId: string
	try {
		const userId = await login(email)
		sessionId = await startSession(userId)
	} catch (err) {
		throw new Error(`could not start new session for user ${email}: ${err}`)
	}
	return sessionId
}

export async function getSession(sessionId: string): Promise<SessionInfo> {
	try {
		const response = await fetch(`${API_BASE_URL}/session/${sessionId}`)
		if (!response.ok) {
			if (response.status === 404) {
				return {}
			}
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error fetching sessions ${response.status} - ${JSON.stringify(errorBody) || 'Unknown error'}`
			)
		}
		const session: SessionInfo = await response.json()
		return session
	} catch (error) {
		console.log('Error fetching problems', error)
		throw error
	}
}

async function login(email: string): Promise<number> {
	try {
		const response = await fetch(`${API_BASE_URL}/login`, {
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

async function startSession(userId: number): Promise<string> {
	try {
		const response = await fetch(`${API_BASE_URL}/session`, {
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
