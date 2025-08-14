import { PUBLIC_API_BASE_URL as API_BASE_URL } from '$env/static/public'

export interface User {
	id: string
	email?: string
	name?: string
	profilePic?: string
}

export interface SessionInfo {
	user?: User
}

export async function createSessionForUser(user: User): Promise<string> {
	let sessionId: string
	try {
		const existingUser = await tryGetExistingUser(user.id)
		if (!existingUser) {
			await createNewUser(user)
		}
		sessionId = await startSession(user.id)
	} catch (err) {
		throw new Error(`could not start new session for user with id ${user.id}: ${err}`)
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

async function tryGetExistingUser(userId: string): Promise<User | undefined> {
	try {
		const response = await fetch(`${API_BASE_URL}/user/${userId}`)
		if (!response.ok) {
			if (response.status === 404) {
				return undefined
			}
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(
				`Error getting existing user ${response.status} - ${errorBody || 'Unknown error'}`
			)
		}
		const user: User = await response.json()
		return user
	} catch (err) {
		throw new Error(`Could not get existing user: ${err}`)
	}
}

async function createNewUser(user: User): Promise<void> {
	try {
		const response = await fetch(`${API_BASE_URL}/user`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(user)
		})
		if (!response.ok) {
			const errorBody = await response.json().catch(() => ({ message: response.statusText }))
			throw new Error(`Error creating user ${response.status} - ${errorBody || 'Unknown error'}`)
		}
	} catch (err) {
		throw new Error(`Could not create new user: ${err}`)
	}
}

async function startSession(userId: string): Promise<string> {
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
