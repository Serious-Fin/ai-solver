package validator

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"serious-fin/api/common"
	"strconv"
	"strings"
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
		return nil, err
	}

	dirPath, err := os.MkdirTemp(".", "test_run_")
	if err != nil {
		return nil, err
	}

	err = createTestFile(fmt.Sprintf("%s/code_test.go", dirPath), body.Code, *testParams)
	if err != nil {
		return nil, err
	}

	testOutput, err := runTests(dirPath)
	if err != nil {
		return nil, err
	}

	testStates, err := parseCommandOutput(testOutput)
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error removing test dir: %v", err)
	}
	return testStates, nil
}

// TODO: Write tests for API
func (vh *ValidatorHandler) fetchTestCreationParams(problemId int) (*testCreationParams, error) {
	var testParams testCreationParams
	row := vh.DB.QueryRow("SELECT testTemplate, testHelpers FROM goTemplates WHERE problemFk = ?", problemId)
	err := row.Scan(&testParams.singleTestTemplate, &testParams.additionalHelpers)
	if err != nil {
		return nil, err
	}

	var testCasesString string
	row = vh.DB.QueryRow("SELECT testCases FROM problems WHERE id = ?", problemId)
	err = row.Scan(&testCasesString)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(testCasesString), &testParams.problemTestCases)
	if err != nil {
		return nil, err
	}
	return &testParams, nil
}

func createTestFile(filename, testableCode string, testParams testCreationParams) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n%s\n", fileStartTemplate, testableCode))
	if err != nil {
		return err
	}

	var newTestCase string
	for _, testCaseData := range testParams.problemTestCases {
		newTestCase = testParams.singleTestTemplate
		newTestCase = strings.Replace(newTestCase, "{{ID}}", fmt.Sprintf("_%s", strconv.Itoa(testCaseData.Id)), 1)
		newTestCase = strings.Replace(newTestCase, "{{OUTPUT}}", testCaseData.ExpectedOutput, 1)
		for inputIndex, input := range testCaseData.Inputs {
			newTestCase = strings.Replace(newTestCase, fmt.Sprintf("{{INPUT%d}}", inputIndex), input, 1)
		}
		_, err = file.WriteString(fmt.Sprintf("%s\n", newTestCase))
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString(testParams.additionalHelpers)
	if err != nil {
		return err
	}
	return nil
}

func runTests(testFilePath string) (string, error) {
	var outputBuffer bytes.Buffer
	testCommand := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", testFilePath), "--network", "none", "go-testing-image:latest", "/bin/sh", "-c", "go mod init test_proj && go test -v")
	testCommand.Stdout = &outputBuffer
	testCommand.Stderr = &outputBuffer

	err := testCommand.Run()
	output := outputBuffer.String()
	if err != nil {
		// return error only if it's status code is other than 1, because failing go tests return exit code 1
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			return output, fmt.Errorf("command execution returned error not of type ExitError: %v", err)
		}
		if exitError.ExitCode() != 1 {
			return output, fmt.Errorf("command execution returned error: %v", err)
		}
	}
	return output, nil
}

func parseCommandOutput(cmdOutput string) (*Response, error) {
	response := &Response{
		SucceededTests: []int{},
		FailedTests:    make([]FailInfo, 0),
	}

	currentTestId := -1
	scanner := bufio.NewScanner(strings.NewReader((cmdOutput)))

	runRegex := regexp.MustCompile(`^=== RUN\s+Test.*_(\d+)$`)
	passRegex := regexp.MustCompile(`^--- PASS:`)
	failRegex := regexp.MustCompile(`^\s+.*?:\d+:\s+got\s+(.*),\s+want\s(.*)$`)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := runRegex.FindStringSubmatch(line); len(matches) > 1 {
			id, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("could not parse test ID from line %s", line)
			}
			currentTestId = id
		} else if passRegex.MatchString(line) && currentTestId != -1 {
			response.SucceededTests = append(response.SucceededTests, currentTestId)
			currentTestId = -1
		} else if matches := failRegex.FindStringSubmatch(line); len(matches) > 2 && currentTestId != -1 {
			response.FailedTests = append(response.FailedTests, FailInfo{
				Id:      currentTestId,
				Got:     matches[1],
				Want:    matches[2],
				Message: WRONG_OUTPUT,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning command output: %w", err)
	}

	return response, nil
}

// TODO: redo validation step using `go test --json`
