import type { TestRunOutput } from '$lib/TestStatusReporter';

const BASE_URL = "http://127.0.0.1:8080"

export interface ValidateRequest {
    problemId: string;
    code: string;
    language: string;
}

export async function validate(req: ValidateRequest): Promise<TestRunOutput> {
    try {
        const resp = await fetch(`${BASE_URL}/validate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(req)
        });
        if (!resp.ok) {
            throw new Error(`HTTP error with status ${resp.status}`);
        }
        const testRunOutput: TestRunOutput = await resp.json();
        return testRunOutput
    } catch (err) {
        throw Error(`Could not call validate endpoint: ${JSON.stringify(err)}`)
    }

}