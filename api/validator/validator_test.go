package validator

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"serious-fin/api/common"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestTextBeforeOutput(t *testing.T) {
	cmdOutput := `
	go: creating new go.mod: module test_proj
	go: to add module requirements and sums:
    go mod tidy
	{"Time":"2025-07-27T18:04:10.225107394Z","Action":"start","Package":"test_proj"}
	{"Time":"2025-07-27T18:04:10.226345745Z","Action":"run","Package":"test_proj","Test":"TestTwoSum_0"}
	{"Time":"2025-07-27T18:04:10.226368636Z","Action":"output","Package":"test_proj","Test":"TestTwoSum_0","Output":"=== RUN   TestTwoSum_0\n"}
	{"Time":"2025-07-27T18:04:10.2263954Z","Action":"output","Package":"test_proj","Test":"TestTwoSum_0","Output":"--- PASS: TestTwoSum_0 (0.00s)\n"}
	{"Time":"2025-07-27T18:04:10.226398054Z","Action":"pass","Package":"test_proj","Test":"TestTwoSum_0","Elapsed":0}`

	want := &Response{
		SucceededTests: []int{0},
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

func TestFetchingCreationParamsBadFirstQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	problemId := 1

	var mockHandler = NewValidatorHandler(db)

	mock.ExpectQuery("SELECT testTemplate, testHelpers FROM goTemplates WHERE problemFk = ?").WithArgs(problemId).WillReturnError(errors.New("error querying data"))

	if _, err = mockHandler.fetchTestCreationParams(problemId); err == nil {
		t.Error("expected error when query fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchingCreationParamsBadSecondQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	problemId := 1

	var mockHandler = NewValidatorHandler(db)

	mock.ExpectQuery("SELECT testTemplate, testHelpers FROM goTemplates WHERE problemFk = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"testTemplate", "testTemplate",
	}).AddRows([]driver.Value{"foo", "bar"}))
	mock.ExpectQuery("SELECT testCases FROM problems WHERE id = ?").WithArgs(problemId).WillReturnError(errors.New("error querying data"))

	if _, err = mockHandler.fetchTestCreationParams(problemId); err == nil {
		t.Error("expected error when query fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchingCreationParamsSecondQueryWrongFormat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	problemId := 1

	var mockHandler = NewValidatorHandler(db)

	mock.ExpectQuery("SELECT testTemplate, testHelpers FROM goTemplates WHERE problemFk = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"testTemplate", "testTemplate",
	}).AddRows([]driver.Value{"foo", "bar"}))
	mock.ExpectQuery("SELECT testCases FROM problems WHERE id = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"testCases",
	}).AddRows([]driver.Value{"bad format"}))

	if _, err = mockHandler.fetchTestCreationParams(problemId); err == nil {
		t.Error("expected error when query fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchingCreationParamsSecondQueryCorrectFormat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	problemId := 1

	var mockHandler = NewValidatorHandler(db)

	mock.ExpectQuery("SELECT testTemplate, testHelpers FROM goTemplates WHERE problemFk = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"testTemplate", "testTemplate",
	}).AddRows([]driver.Value{"foo", "bar"}))
	mock.ExpectQuery("SELECT testCases FROM problems WHERE id = ?").WithArgs(problemId).WillReturnRows(sqlmock.NewRows([]string{
		"testCases",
	}).AddRows([]driver.Value{`[{"id": 0,"inputs":  ["[]int{2, 7, 11, 15}","9"],"output": "[]int{0, 1}"}]`}))

	if _, err = mockHandler.fetchTestCreationParams(problemId); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateTestFileNonExistentPath(t *testing.T) {
	dirPath := "/NON/EXISTENT/PATH/code.go"
	err := createTestFile(dirPath, "code", testCreationParams{})
	if err == nil {
		t.Errorf("expected error when path \"%s\"does not exist", dirPath)
	}

	if _, ok := err.(*os.PathError); !ok {
		t.Error("expected error to be PathError")

		// delete file because it was created when it shouldn't have been
		err := os.RemoveAll(dirPath)
		if err != nil {
			t.Errorf("could not remove files that should not have been created. remove files at path \"%s\" manually", dirPath)
		}
	}
}

func TestCreateTestFileIsStartTemplateAdded(t *testing.T) {
	dirPath := fmt.Sprintf("./test_file_%s", randSeq(4))
	err := createTestFile(dirPath, "code", testCreationParams{})
	if err != nil {
		t.Errorf("unexpected error when creating file \"%s\": %v", dirPath, err)
	}

	content, err := os.ReadFile(dirPath)
	if err != nil {
		t.Errorf("failed to read created file: %v", err)
	}
	fileContents := string(content)
	if !strings.HasPrefix(fileContents, fileStartTemplate) {
		t.Errorf("File content does not start with expected template.\nExpected to find:\n%s\nFile contents:\n%s", fileStartTemplate, fileContents)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		t.Errorf("could not remove files that should not have been created. remove files at path \"%s\" manually", dirPath)
	}
}

func TestCreateTestFileIsUserCodeAdded(t *testing.T) {
	dirPath := fmt.Sprintf("./test_file_%s", randSeq(4))
	userCode := "foo bar baz"
	err := createTestFile(dirPath, userCode, testCreationParams{})
	if err != nil {
		t.Errorf("unexpected error when creating file \"%s\": %v", dirPath, err)
	}

	content, err := os.ReadFile(dirPath)
	if err != nil {
		t.Errorf("failed to read created file: %v", err)
	}
	fileContents := string(content)
	if !strings.Contains(fileContents, userCode) {
		t.Errorf("File content does not start with expected template.\nExpected to find:\n%s\nFile contents:\n%s", userCode, fileContents)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		t.Errorf("could not remove files that should not have been created. remove files at path \"%s\" manually", dirPath)
	}
}

func TestCreateTestFileIsTestCodeAdded(t *testing.T) {
	dirPath := fmt.Sprintf("./test_file_%s", randSeq(4))
	testCases := []common.TestCase{
		{
			Id:             0,
			Inputs:         []string{"[]string{1, 2, 3}", "\"foo bar baz\""},
			ExpectedOutput: "15.5",
		},
		{
			Id:             1,
			Inputs:         []string{"[]string{4, 5, 6}", "\"what is this\""},
			ExpectedOutput: "88.8",
		},
	}
	testTemplate := `func testing{{ID}}(t * testing.T) {
	want := {{OUTPUT}}
	got := runFunc({{INPUT0}}, {{INPUT1}})
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}`
	want := `func testing_0(t * testing.T) {
	want := 15.5
	got := runFunc([]string{1, 2, 3}, "foo bar baz")
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
func testing_1(t * testing.T) {
	want := 88.8
	got := runFunc([]string{4, 5, 6}, "what is this")
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}`
	err := createTestFile(dirPath, "code", testCreationParams{
		singleTestTemplate: testTemplate,
		problemTestCases:   testCases,
	})
	if err != nil {
		t.Errorf("unexpected error when creating file \"%s\": %v", dirPath, err)
	}

	content, err := os.ReadFile(dirPath)
	if err != nil {
		t.Errorf("failed to read created file: %v", err)
	}
	fileContents := string(content)
	if !strings.Contains(fileContents, want) {
		t.Errorf("File content does not start with expected template.\nExpected to find:\n%s\nFile contents:\n%s", want, fileContents)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		t.Errorf("could not remove files that should not have been created. remove files at path \"%s\" manually", dirPath)
	}
}

func TestCreateTestFileIsHelperCodeAdded(t *testing.T) {
	dirPath := fmt.Sprintf("./test_file_%s", randSeq(4))
	helpers := "helper functions"
	err := createTestFile(dirPath, "code", testCreationParams{
		additionalHelpers: helpers,
	})
	if err != nil {
		t.Errorf("unexpected error when creating file \"%s\": %v", dirPath, err)
	}

	content, err := os.ReadFile(dirPath)
	if err != nil {
		t.Errorf("failed to read created file: %v", err)
	}
	fileContents := string(content)
	if !strings.Contains(fileContents, helpers) {
		t.Errorf("File content does not start with expected template.\nExpected to find:\n%s\nFile contents:\n%s", helpers, fileContents)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		t.Errorf("could not remove files that should not have been created. remove files at path \"%s\" manually", dirPath)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
