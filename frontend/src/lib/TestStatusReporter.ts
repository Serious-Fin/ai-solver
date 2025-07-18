enum TestStatus{
    UNKNOWN, PASS, FAIL
}

interface SingleTestStatus {
    status: TestStatus,
    want?: string,
    got?: string
}

interface TestRunOutput {
    succeededTests: number[];
    failedTests: {[id: number]: TestFailReason}
}

interface TestFailReason {
    want: string,
    got: string,
    message: string
}

export class TestStatusReporter {
    private statuses: { [id: number]: SingleTestStatus } = {};

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

        output.failedTests.forEach(id => {
            this.statuses[id].status = TestStatus.FAIL
            this.statuses[id].want = 
        })
    }
}

// TODO: failed test output better be array of objects instead of hash