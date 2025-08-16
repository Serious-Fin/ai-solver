package validator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"serious-fin/api/common"
	"strconv"
	"strings"
	"time"
)

type ValidatorHandler struct {
	DB common.DBInterface
}

type Request struct {
	ProblemId int    `form:"problemId"`
	Code      string `form:"code"`
}

type Response struct {
	FailedTests    []FailInfo `json:"failedTests"`
	SucceededTests []int      `json:"succeededTests"`
}

type FailInfo struct {
	Id      int    `json:"id"`
	Want    string `json:"want"`
	Got     string `json:"got"`
	Message string `json:"message"`
}

type testEvent struct {
	Time        time.Time `json:"Time"`
	Action      string    `json:"Action"`
	Package     string    `json:"Package"`
	Test        string    `json:"Test,omitempty"`
	Elapsed     float64   `json:"Elapsed,omitempty"`
	Output      string    `json:"Output,omitempty"`
	FailedBuild string    `json:"FailedBuild,omitempty"`
}

type testCreationParams struct {
	singleTestTemplate string
	additionalHelpers  string
	problemTestCases   []common.TestCase
}

var fileStartTemplate = `package main
import "testing"
`

const (
	WRONG_OUTPUT = "wrong output"
)

func NewValidatorHandler(db common.DBInterface) *ValidatorHandler {
	return &ValidatorHandler{
		DB: db,
	}
}

func (vh *ValidatorHandler) Validate(body Request) (*Response, error) {
	testParams, err := vh.fetchTestCreationParams(body.ProblemId)
	if err != nil {
		return nil, fmt.Errorf("could not fetch test creation params: %w", err)
	}

	dirPath, err := os.MkdirTemp(".", "test_run_")
	if err != nil {
		return nil, fmt.Errorf("error making temporary directory: %w", err)
	}
	defer os.RemoveAll(dirPath)

	err = createTestFile(fmt.Sprintf("%s/code_test.go", dirPath), body.Code, *testParams)
	if err != nil {
		return nil, fmt.Errorf("error creating test file: %w", err)
	}

	testOutput, err := runTests(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error running tests in docker: %w", err)
	}

	testStates, err := parseCommandOutput(testOutput)
	if err != nil {
		return nil, fmt.Errorf("error parsing command output %s: %w", testOutput, err)
	}

	return testStates, nil
}

func (vh *ValidatorHandler) fetchTestCreationParams(problemId int) (*testCreationParams, error) {
	var testParams testCreationParams
	row := vh.DB.QueryRow("SELECT testTemplate, testHelpers FROM goTemplates WHERE problemFk = ?", problemId)
	err := row.Scan(&testParams.singleTestTemplate, &testParams.additionalHelpers)
	if err != nil {
		return nil, fmt.Errorf("error scanning templates and helpers from db (problem id %d): %w", problemId, err)
	}

	var testCasesString string
	row = vh.DB.QueryRow("SELECT testCases FROM problems WHERE id = ?", problemId)
	err = row.Scan(&testCasesString)
	if err != nil {
		return nil, fmt.Errorf("error scanning test cases from db (problem id %d): %w", problemId, err)
	}

	err = json.Unmarshal([]byte(testCasesString), &testParams.problemTestCases)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal test cases from string \"%s\" (problem id %d): %w", testCasesString, problemId, err)
	}
	return &testParams, nil
}

func createTestFile(filename, testableCode string, testParams testCreationParams) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file with name \"%s\", error: %w", filename, err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n%s\n", fileStartTemplate, testableCode)
	if err != nil {
		return fmt.Errorf("could not write start template and user code to file: %w", err)
	}

	var newTestCase string
	for _, testCaseData := range testParams.problemTestCases {
		newTestCase = testParams.singleTestTemplate
		newTestCase = strings.Replace(newTestCase, "{{ID}}", fmt.Sprintf("_%s", strconv.Itoa(testCaseData.Id)), 1)
		newTestCase = strings.Replace(newTestCase, "{{OUTPUT}}", testCaseData.ExpectedOutput, 1)
		for inputIndex, input := range testCaseData.Inputs {
			newTestCase = strings.Replace(newTestCase, fmt.Sprintf("{{INPUT%d}}", inputIndex), input, 1)
		}
		_, err = fmt.Fprintf(file, "%s\n", newTestCase)
		if err != nil {
			return fmt.Errorf("could not write test case to file: %w", err)
		}
	}

	_, err = file.WriteString(testParams.additionalHelpers)
	if err != nil {
		return fmt.Errorf("could not write additional helper functions to file: %w", err)
	}
	return nil
}

func runTests(testFilePath string) (string, error) {
	initCmd := exec.Command("go", "mod", "init", "test_proj")
	initCmd.Dir = testFilePath
	if err := initCmd.Run(); err != nil {
		return "", fmt.Errorf("go mod init failed: %w", err)
	}

	testCmd := exec.Command("go", "test", "-json")
	testCmd.Dir = testFilePath
	output, err := testCmd.Output()
	if err != nil {
		// return error only if it's status code is other than 1, because failing go tests return exit code 1
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			return string(output), fmt.Errorf("command execution returned error not of type ExitError: %w", err)
		}
		if exitError.ExitCode() != 1 {
			return string(output), fmt.Errorf("command execution returned: %s, error: %w", string(output), err)
		}
	}
	return string(output), nil
}

func parseCommandOutput(cmdOutput string) (*Response, error) {
	response := &Response{
		SucceededTests: []int{},
		FailedTests:    make([]FailInfo, 0),
	}
	testOutputs := make(map[int][]string)

	cmdOutput = strings.TrimSpace(cmdOutput)
	scanner := bufio.NewScanner(strings.NewReader((cmdOutput)))
	var line string
	for scanner.Scan() {
		line = scanner.Text()

		var testLog testEvent
		_ = json.Unmarshal([]byte(line), &testLog)

		if testLog.Test == "" {
			// test event log is not associated with any specific test so we skip this log
			continue
		}

		switch testLog.Action {
		case "output":
			testId, err := getTestId(testLog.Test)
			if err != nil {
				return nil, fmt.Errorf("could not get test id from output event: %w", err)
			}
			testOutputs[testId] = append(testOutputs[testId], testLog.Output)
		case "pass":
			testId, err := getTestId(testLog.Test)
			if err != nil {
				return nil, fmt.Errorf("could not get test id from pass event: %w", err)
			}
			delete(testOutputs, testId)
			response.SucceededTests = append(response.SucceededTests, testId)
		case "fail":
			testId, err := getTestId(testLog.Test)
			if err != nil {
				return nil, fmt.Errorf("could not get test id from fail event: %w", err)
			}

			foundGotAndWant := false
			for _, output := range testOutputs[testId] {
				got, want, err := getGotWantValues(output)
				if nil == err {
					response.FailedTests = append(response.FailedTests, FailInfo{
						Id:      testId,
						Want:    want,
						Got:     got,
						Message: WRONG_OUTPUT,
					})
					delete(testOutputs, testId)
					foundGotAndWant = true
					break
				}
			}
			if foundGotAndWant {
				continue
			}

			return nil, fmt.Errorf("did not find \"got\" and \"want\" values for failed test \"%s\" in output values %v", testLog.Test, testOutputs[testId])
		}
	}

	return response, nil
}

var testIdFromNameRegex = regexp.MustCompile(`.+_(\d+)$`)

func getTestId(testName string) (int, error) {
	matches := testIdFromNameRegex.FindStringSubmatch(testName)
	if len(matches) != 2 {
		return -1, fmt.Errorf("could not find one test id match in test name \"%s\", found matches %d (expected 2)", testName, len(matches))
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, fmt.Errorf("could not parse test id from match \"%s\"", matches[1])
	}
	return id, nil
}

var gotAndWantValueRegex = regexp.MustCompile(`got\s+(.*),\s+want\s+(.*)\n$`)

func getGotWantValues(text string) (string, string, error) {
	matches := gotAndWantValueRegex.FindStringSubmatch(text)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("could not find got and want values in \"%s\", found matches %d (expected 3)", text, len(matches))
	}
	return matches[1], matches[2], nil
}
