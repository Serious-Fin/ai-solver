package validator

import (
	"reflect"
	"testing"
)

func TestValidationOutputParsing(t *testing.T) {
	cmdOutput := `go: creating new go.mod: module test_proj
go: to add module requirements and sums:
        go mod tidy
=== RUN   TestTwoSum_0
--- PASS: TestTwoSum_0 (0.00s)
=== RUN   TestTwoSum_1
--- PASS: TestTwoSum_1 (0.00s)
=== RUN   TestTwoSum_2
--- PASS: TestTwoSum_2 (0.00s)
=== RUN   TestTwoSum_3
--- PASS: TestTwoSum_3 (0.00s)
=== RUN   TestTwoSum_4
    code_test.go:48: got [], want [4 9]
--- FAIL: TestTwoSum_4 (0.00s)
=== RUN   TestTwoSum_5
    code_test.go:55: got [1 3], want [1 4]
--- FAIL: TestTwoSum_5 (0.00s)
=== RUN   TestTwoSum_6
--- PASS: TestTwoSum_6 (0.00s)
=== RUN   TestTwoSum_7
--- PASS: TestTwoSum_7 (0.00s)
=== RUN   TestTwoSum_8
    code_test.go:76: got [4 5], want [3 7]
--- FAIL: TestTwoSum_8 (0.00s)
=== RUN   TestTwoSum_9
--- PASS: TestTwoSum_9 (0.00s)
FAIL
exit status 1
FAIL    test_proj       0.002s
`
	want := &Response{
		SucceededTests: []int{0, 1, 2, 3, 6, 7, 9},
		FailedTests: []FailInfo{
			{
				Id:      4,
				Want:    "[4 9]",
				Got:     "[]",
				Message: "wrong output",
			},
			{
				Id:      5,
				Want:    "[1 4]",
				Got:     "[1 3]",
				Message: "wrong output",
			},
			{
				Id:      8,
				Want:    "[3 7]",
				Got:     "[4 5]",
				Message: "wrong output",
			},
		},
	}

	got, err := ParseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAllPassingCase(t *testing.T) {
	cmdOutput := `=== RUN   TestTwoSum_0
--- PASS: TestTwoSum_0 (0.00s)
=== RUN   TestTwoSum_1
--- PASS: TestTwoSum_1 (0.00s)
=== RUN   TestTwoSum_2
--- PASS: TestTwoSum_2 (0.00s)
=== RUN   TestTwoSum_3
--- PASS: TestTwoSum_3 (0.00s)
=== RUN   TestTwoSum_4
--- PASS: TestTwoSum_4 (0.00s)
=== RUN   TestTwoSum_5
--- PASS: TestTwoSum_5 (0.00s)
=== RUN   TestTwoSum_6
--- PASS: TestTwoSum_6 (0.00s)
=== RUN   TestTwoSum_7
--- PASS: TestTwoSum_7 (0.00s)
=== RUN   TestTwoSum_8
--- PASS: TestTwoSum_8 (0.00s)
=== RUN   TestTwoSum_9
--- PASS: TestTwoSum_9 (0.00s)
PASS
exit status 0
PASS    test_proj       0.002s
`
	want := &Response{
		SucceededTests: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		FailedTests:    []FailInfo{},
	}

	got, err := ParseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAllFailingCase(t *testing.T) {
	cmdOutput := `=== RUN   TestTwoSum_0
    code_test.go:48: got hello world, want goodbye world
--- FAIL: TestTwoSum_0 (0.00s)
=== RUN   TestTwoSum_1
    code_test.go:48: got foo bar baz, want apple banana
--- FAIL: TestTwoSum_1 (0.00s)
FAIL
exit status 1
FAIL    test_proj       0.002s
`
	want := &Response{
		SucceededTests: []int{},
		FailedTests: []FailInfo{
			{
				Id:      0,
				Want:    "goodbye world",
				Got:     "hello world",
				Message: "wrong output",
			},
			{
				Id:      1,
				Want:    "apple banana",
				Got:     "foo bar baz",
				Message: "wrong output",
			},
		},
	}

	got, err := ParseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

var saved = `
{"Time":"2025-07-23T17:39:44.806464+03:00","Action":"start","Package":"serious-fin/api/problem"}
{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestGetProblemsQueryThrowsError"}
{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetProblemsQueryThrowsError","Output":"=== RUN   TestGetProblemsQueryThrowsError\n"}
{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetProblemsQueryThrowsError","Output":"--- PASS: TestGetProblemsQueryThrowsError (0.00s)\n"}
{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestGetProblemsQueryThrowsError","Elapsed":0}
{"Time":"2025-07-23T17:39:44.977998+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestGetProblems"}
{"Time":"2025-07-23T17:39:44.978004+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetProblems","Output":"=== RUN   TestGetProblems\n"}
{"Time":"2025-07-23T17:39:44.978091+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetProblems","Output":"--- PASS: TestGetProblems (0.00s)\n"}
{"Time":"2025-07-23T17:39:44.978097+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestGetProblems","Elapsed":0}
{"Time":"2025-07-23T17:39:44.978103+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestGetProblemById"}
{"Time":"2025-07-23T17:39:44.978107+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetProblemById","Output":"=== RUN   TestGetProblemById\n"}
{"Time":"2025-07-23T17:39:44.978294+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetProblemById","Output":"--- PASS: TestGetProblemById (0.00s)\n"}
{"Time":"2025-07-23T17:39:44.978304+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestGetProblemById","Elapsed":0}
{"Time":"2025-07-23T17:39:44.978309+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestGetMainFuncGo"}
{"Time":"2025-07-23T17:39:44.978312+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetMainFuncGo","Output":"=== RUN   TestGetMainFuncGo\n"}
{"Time":"2025-07-23T17:39:44.978332+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestGetMainFuncGo","Output":"--- PASS: TestGetMainFuncGo (0.00s)\n"}
{"Time":"2025-07-23T17:39:44.978336+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestGetMainFuncGo","Elapsed":0}
{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestPurposefulBad1"}
{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestPurposefulBad1","Output":"=== RUN   TestPurposefulBad1\n"}
{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestPurposefulBad1","Output":"    problem_test.go:161: got [foo bar], want [apple banana]\n"}
{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestPurposefulBad1","Output":"--- FAIL: TestPurposefulBad1 (0.00s)\n"}
{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestPurposefulBad1","Elapsed":0}
{"Time":"2025-07-23T17:39:44.978384+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\n"}
{"Time":"2025-07-23T17:39:44.978865+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"exit status 1\n"}
{"Time":"2025-07-23T17:39:44.978872+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\tserious-fin/api/problem\t0.172s\n"}
{"Time":"2025-07-23T17:39:44.978877+03:00","Action":"fail","Package":"serious-fin/api/problem","Elapsed":0.172}`
