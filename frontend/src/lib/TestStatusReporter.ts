import type { TestCase } from './api/problems'

export enum TestStatus {
	UNKNOWN,
	PASS,
	FAIL
}

export interface SingleTestStatus {
	id: number
	status: TestStatus
	inputs: string[]
	output: string
	got?: string
}

export interface TestRunOutput {
	succeededTests: number[]
	failedTests: FailReason[]
}

interface FailReason {
	id: number
	got: string
	message: string
}

export class TestStatusReporter {
	private statuses: { [id: number]: SingleTestStatus } = {}

	get getTestStatuses() {
		return this.statuses
	}

	public constructor(testCases: TestCase[]) {
		testCases.forEach((testCase) => {
			this.statuses[testCase.id] = {
				id: testCase.id,
				status: TestStatus.UNKNOWN,
				inputs: testCase.inputs,
				output: testCase.output
			}
		})
	}

	public UpdateTestStatuses(output: TestRunOutput) {
		output.succeededTests.forEach((id) => {
			this.statuses[id].status = TestStatus.PASS
		})

		output.failedTests.forEach((failReason) => {
			this.statuses[failReason.id].status = TestStatus.FAIL
			this.statuses[failReason.id].got = failReason.got
		})
	}

	public GetTestStatuses(): SingleTestStatus[] {
		const testResults: SingleTestStatus[] = []
		for (const key in this.statuses) {
			testResults.push(this.statuses[key])
		}
		return testResults
	}

	public AllTestsSuccessful(): boolean {
		for (const key in this.statuses) {
			if (
				this.statuses[key].status === TestStatus.FAIL ||
				this.statuses[key].status === TestStatus.UNKNOWN
			) {
				return false
			}
		}
		return true
	}
}
