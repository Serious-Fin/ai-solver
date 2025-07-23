package validator

import (
	"reflect"
	"testing"
)

func TestValidationOutputParsing(t *testing.T) {
	cmdOutput := `
	{"Time":"2025-07-23T17:39:44.806464+03:00","Action":"start","Package":"serious-fin/api/problem"}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_0"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"=== RUN   TestTwoSum_0\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"--- PASS: TestTwoSum_0 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_1"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"=== RUN   TestTwoSum_1\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"--- PASS: TestTwoSum_1 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_2"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"=== RUN   TestTwoSum_2\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"--- PASS: TestTwoSum_2 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_3"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_3","Output":"=== RUN   TestTwoSum_3\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_3","Output":"--- PASS: TestTwoSum_3 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_3","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_4"}
	{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_4","Output":"=== RUN   TestTwoSum_4\n"}
	{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_4","Output":"    code_test.go:48: got [], want [4 9]\n"}
	{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_4","Output":"--- FAIL: TestTwoSum_4 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestTwoSum_4","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_5"}
	{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_5","Output":"=== RUN   TestTwoSum_5\n"}
	{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_5","Output":"    code_test.go:55: got [1 3], want [1 4]\n"}
	{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_5","Output":"--- FAIL: TestTwoSum_5 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestTwoSum_5","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_6"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_6","Output":"=== RUN   TestTwoSum_6\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_6","Output":"--- PASS: TestTwoSum_6 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_6","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_7"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_7","Output":"=== RUN   TestTwoSum_7\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_7","Output":"--- PASS: TestTwoSum_7 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_7","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_8"}
	{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_8","Output":"=== RUN   TestTwoSum_8\n"}
	{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_8","Output":"    code_test.go:76: got [4 5], want [3 7]\n"}
	{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_8","Output":"--- FAIL: TestTwoSum_8 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestTwoSum_8","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_9"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_9","Output":"=== RUN   TestTwoSum_9\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_9","Output":"--- PASS: TestTwoSum_9 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_9","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.978384+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\n"}
	{"Time":"2025-07-23T17:39:44.978865+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"exit status 1\n"}
	{"Time":"2025-07-23T17:39:44.978872+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\ttest_proj\t0.002s\n"}
	{"Time":"2025-07-23T17:39:44.978877+03:00","Action":"fail","Package":"serious-fin/api/problem","Elapsed":0.172}`

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

	got, err := parseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAllPassingCase(t *testing.T) {
	cmdOutput := `
	{"Time":"2025-07-23T17:39:44.806464+03:00","Action":"start","Package":"serious-fin/api/problem"}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_0"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"=== RUN   TestTwoSum_0\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"--- PASS: TestTwoSum_0 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_1"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"=== RUN   TestTwoSum_1\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"--- PASS: TestTwoSum_1 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_2"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"=== RUN   TestTwoSum_2\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"--- PASS: TestTwoSum_2 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_3"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_3","Output":"=== RUN   TestTwoSum_3\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_3","Output":"--- PASS: TestTwoSum_3 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_3","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.978384+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"PASS\n"}
	{"Time":"2025-07-23T17:39:44.978865+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"exit status 0\n"}
	{"Time":"2025-07-23T17:39:44.978872+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"PASS\ttest_proj\t0.002s\n"}
	{"Time":"2025-07-23T17:39:44.978877+03:00","Action":"fail","Package":"serious-fin/api/problem","Elapsed":0.172}`

	want := &Response{
		SucceededTests: []int{0, 1, 2, 3},
		FailedTests:    []FailInfo{},
	}

	got, err := parseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAllFailingCase(t *testing.T) {
	cmdOutput := `
	{"Time":"2025-07-23T17:39:44.806464+03:00","Action":"start","Package":"serious-fin/api/problem"}
	{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_0"}
	{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"=== RUN   TestTwoSum_0\n"}
	{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"    code_test.go:48: got hello world, want goodbye world\n"}
	{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"--- FAIL: TestTwoSum_0 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_1"}
	{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"=== RUN   TestTwoSum_1\n"}
	{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"    code_test.go:48: got foo bar baz, want apple banana\n"}
	{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"--- FAIL: TestTwoSum_1 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.978384+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\n"}
	{"Time":"2025-07-23T17:39:44.978865+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"exit status 1\n"}
	{"Time":"2025-07-23T17:39:44.978872+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\ttest_proj\t0.002s\n"}
	{"Time":"2025-07-23T17:39:44.978877+03:00","Action":"fail","Package":"serious-fin/api/problem","Elapsed":0.172}`

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

	got, err := parseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSuccessAndFailureWithArrays(t *testing.T) {
	cmdOutput := `
	{"Time":"2025-07-23T17:39:44.806464+03:00","Action":"start","Package":"serious-fin/api/problem"}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_0"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"=== RUN   TestTwoSum_0\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Output":"--- PASS: TestTwoSum_0 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_0","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.977589+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_1"}
	{"Time":"2025-07-23T17:39:44.977634+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"=== RUN   TestTwoSum_1\n"}
	{"Time":"2025-07-23T17:39:44.977976+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Output":"--- PASS: TestTwoSum_1 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.977987+03:00","Action":"pass","Package":"serious-fin/api/problem","Test":"TestTwoSum_1","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.97834+03:00","Action":"run","Package":"serious-fin/api/problem","Test":"TestTwoSum_2"}
	{"Time":"2025-07-23T17:39:44.978342+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"=== RUN   TestTwoSum_2\n"}
	{"Time":"2025-07-23T17:39:44.978362+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"    problem_test.go:161: got [foo bar], want [apple banana]\n"}
	{"Time":"2025-07-23T17:39:44.978374+03:00","Action":"output","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Output":"--- FAIL: TestTwoSum_2 (0.00s)\n"}
	{"Time":"2025-07-23T17:39:44.97838+03:00","Action":"fail","Package":"serious-fin/api/problem","Test":"TestTwoSum_2","Elapsed":0}
	{"Time":"2025-07-23T17:39:44.978384+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\n"}
	{"Time":"2025-07-23T17:39:44.978865+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"exit status 1\n"}
	{"Time":"2025-07-23T17:39:44.978872+03:00","Action":"output","Package":"serious-fin/api/problem","Output":"FAIL\tserious-fin/api/problem\t0.172s\n"}
	{"Time":"2025-07-23T17:39:44.978877+03:00","Action":"fail","Package":"serious-fin/api/problem","Elapsed":0.172}`

	want := &Response{
		SucceededTests: []int{0, 1},
		FailedTests: []FailInfo{
			{
				Id:      2,
				Got:     "[foo bar]",
				Want:    "[apple banana]",
				Message: WRONG_OUTPUT,
			},
		},
	}

	got, err := parseCommandOutput(cmdOutput)
	if err != nil {
		t.Errorf("error while parsing: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
