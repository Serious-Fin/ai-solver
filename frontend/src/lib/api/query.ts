export interface QueryRequest {
    input: string;
    code: string;
    language: string;
    agent: string;
    sessionId: string;
}

export interface QueryResponse {
    response: string
}

const BASE_URL = "http://127.0.0.1:8080"

export async function query(params: QueryRequest): Promise<string> {
    try {
        const response = await fetch(`${BASE_URL}/query/${params.sessionId}`, {
            body: JSON.stringify({
                input: params.input,
                code: params.code,
                language: params.language,
                agent: params.agent
            }),
            headers: {
                "Content-Type": "application/json"
            },
            method: "POST"
        })
        if (!response.ok) {
            const errorBody = await response.json().catch(() => ({ message: response.statusText }))
            throw new Error(`Error querying agent ${response.status} - ${errorBody || "Unknown error"}`)
        }
        const parsedResp: QueryResponse = await response.json()
        return parsedResp.response
    } catch (err) {
        throw Error(`Could not call query endpoint: ${JSON.stringify(err)}`)
    }
}
