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
	Language  string `form:"language"`
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

type TestParams struct {
	Template string
	Helpers  string
	Cases    []common.TestCase
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
	testParams, err := vh.FetchTestDetails(body.Language, body.ProblemId)
	if err != nil {
		return nil, err
	}

	dirPath, err := os.MkdirTemp(".", "test_run_")
	if err != nil {
		return nil, err
	}
	CreateTestFile(fmt.Sprintf("%s/code_test.go", dirPath), body.Code, testParams.Template, testParams.Cases, testParams.Helpers)
	var outputBuffer bytes.Buffer
	testCommand := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", dirPath), "--network", "none", "go-testing-image:latest", "/bin/sh", "-c", "go mod init test_proj && go test -v")
	testCommand.Stdout = &outputBuffer
	testCommand.Stderr = &outputBuffer

	_ = testCommand.Run()
	testStates, err := ParseCommandOutput(outputBuffer.String())
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		fmt.Printf("Could not remove test dir: %v\n", err)
	}

	return testStates, nil
}

// TODO: Write tests for API
func (vh *ValidatorHandler) FetchTestDetails(language string, problemId int) (*TestParams, error) {
	var testParams TestParams
	var testCasesString string
	sqlString := fmt.Sprintf("SELECT testCases, %sTestTemplate, %sTestHelpers FROM problems WHERE id = ?;", language, language)
	row := vh.DB.QueryRow(sqlString, problemId)
	err := row.Scan(&testCasesString, &testParams.Template, &testParams.Helpers)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(testCasesString), &testParams.Cases)
	if err != nil {
		return nil, err
	}
	return &testParams, nil
}

func CreateTestFile(filename, userCode, testTemplate string, testCases []common.TestCase, helperFuncs string) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", fileStartTemplate))
	check(err)
	_, err = file.WriteString(fmt.Sprintf("%s\n", userCode))
	check(err)

	for _, testCase := range testCases {
		newTestCode := testTemplate
		newTestCode = strings.Replace(newTestCode, "{{ID}}", fmt.Sprintf("_%s", strconv.Itoa(testCase.Id)), 1)
		newTestCode = strings.Replace(newTestCode, "{{OUTPUT}}", testCase.ExpectedOutput, 1)
		for inputIndex, input := range testCase.Inputs {
			newTestCode = strings.Replace(newTestCode, fmt.Sprintf("{{INPUT%d}}", inputIndex), input, 1)
		}
		_, err = file.WriteString(fmt.Sprintf("%s\n", newTestCode))
		check(err)
	}

	_, err = file.WriteString(helperFuncs)
	check(err)
}

func ParseCommandOutput(cmdOutput string) (*Response, error) {
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: redo validation step using `go test --json`
