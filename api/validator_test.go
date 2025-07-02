package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
	want := ValidateResponse{
		SucceededTests: []int{0, 1, 2, 3, 6, 7, 9},
		FailedTests: map[int]FailReason{
			4: {
				Want:    "[4 9]",
				Got:     "[]",
				Message: "Wrong output",
			},
			5: {
				Want:    "[1 4]",
				Got:     "[1 3]",
				Message: "Wrong output",
			},
			8: {
				Want:    "[3 7]",
				Got:     "[4 5]",
				Message: "Wrong output",
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

func ParseCommandOutput(cmdOutput string) (*ValidateResponse, error) {
	var response ValidateResponse

	// split everything on === RUN
	slices := strings.SplitSeq(cmdOutput, "=== RUN")
	for slice := range slices {
		// get  test id
		re := regexp.MustCompile(`Test.*_\d+`)
		regexMatches := re.FindStringSubmatch(slice)
		if len(regexMatches) == 0 {
			continue
		}
		testName := regexMatches[0]
		re = regexp.MustCompile(`\d+$`)
		regexMatches = re.FindStringSubmatch(testName)
		if len(regexMatches) == 0 {
			continue
		}
		testId, err := strconv.Atoi(regexMatches[0])
		if err != nil {
			return nil, err
		}
		fmt.Println(testId)

		// check if test succeeded
		re = regexp.MustCompile(`--- PASS`)
		regexMatches = re.FindStringSubmatch(slice)
		if len(regexMatches) > 0 {
			response.SucceededTests = append(response.SucceededTests, testId)
		}

		// check if test failed
		re = regexp.MustCompile(`--- FAIL`)
		regexMatches = re.FindStringSubmatch(testName)
		if len(regexMatches) == 0 {
			continue
		}

		// find the "want" and "got"
		re = regexp.MustCompile()
	}
	return &response, nil
}
