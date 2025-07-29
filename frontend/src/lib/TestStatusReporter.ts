export enum TestStatus {
    UNKNOWN, PASS, FAIL
}

export interface SingleTestStatus {
    id: number;
    status: TestStatus;
    want?: string;
    got?: string;
}

export interface TestRunOutput {
    succeededTests: number[];
    failedTests: FailReason[];
}

interface FailReason {
    id: number
    want: string;
    got: string;
    message: string;
}

export class TestStatusReporter {
    private statuses: { [id: number]: SingleTestStatus } = {};

    get getTestStatuses() {
        return this.statuses
    }

    public constructor(testCaseIds: number[]) {
        testCaseIds.forEach(id => {
            this.statuses[id] = {
                id,
                status: TestStatus.UNKNOWN
            }
        });
    }

    public UpdateTestStatuses(output: TestRunOutput) {
        output.succeededTests.forEach(id => {
            this.statuses[id].status = TestStatus.PASS
        })

        output.failedTests.forEach(failReason => {
            this.statuses[failReason.id].status = TestStatus.FAIL
            this.statuses[failReason.id].want = failReason.want
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
}

// TODO: lint the project
// sort the tests?