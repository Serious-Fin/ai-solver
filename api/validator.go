package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidateRequest struct {
	ProblemId string `form:"problemId"`
	Code      string `form:"code"`
	Language  string `form:"language"`
}

type FailReason struct {
	Want    string `json:"want"`
	Got     string `json:"got"`
	Message string `json:"message"`
}

type ValidateResponse struct {
	FailedTests    map[int]FailReason `json:"failed"`
	SucceededTests []int              `json:"succeeded"`
}

type TestParams struct {
	Template string
	Helpers  string
	Cases    []TestCase
}

var fileStartTemplate = `package main
import "testing"
`

const (
	WRONG_OUTPUT = "wrong output"
)

func validateCode(c *gin.Context) {
	var body ValidateRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	testParams, err := FetchTestDetails(body.Language, body.ProblemId)
	if err != nil {
		c.Error(err)
		return
	}

	dirPath, err := os.MkdirTemp(".", "test_run_")
	if err != nil {
		c.Error(err)
		return
	}
	CreateTestFile(fmt.Sprintf("%s/code_test.go", dirPath), body.Code, testParams.Template, testParams.Cases, testParams.Helpers)
	var outputBuffer bytes.Buffer
	testCommand := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", dirPath), "--network", "none", "go-testing-image:latest", "/bin/sh", "-c", "go mod init test_proj && go test -v")
	testCommand.Stdout = &outputBuffer
	testCommand.Stderr = &outputBuffer

	err = testCommand.Run()
	if err != nil {
		fmt.Println("--- Err Output ---")
		fmt.Println(outputBuffer.String())
		fmt.Println("----------------------")

		c.Error(err)
		return
	}

	fmt.Println(outputBuffer.String())

	// validate
	// run the following:
	//  docker run -it --rm -v ./test_run_2531706648:/app --network none go-testing-image:latest /bin/sh -c "go mod init test_proj && go test -v"
	// remove temp folder
	c.IndentedJSON(http.StatusOK, ValidateResponse{})
}

// TODO: make validation a long running process: first POST request creates a validation request and subsequent GET requests get the status

func FetchTestDetails(language string, problemId string) (*TestParams, error) {
	var testParams TestParams
	var testCasesString string
	sqlString := fmt.Sprintf("SELECT testCases, %sTestTemplate, %sTestHelpers FROM problems WHERE id = ?;", language, language)
	row := db.QueryRow(sqlString, problemId)
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

func CreateTestFile(filename string, userCode string, testTemplate string, testCases []TestCase, helperFuncs string) {
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

func ParseCommandOutput(cmdOutput string) (*ValidateResponse, error) {
	response := &ValidateResponse{
		FailedTests: make(map[int]FailReason),
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
			response.FailedTests[currentTestId] = FailReason{
				Got:     matches[1],
				Want:    matches[2],
				Message: WRONG_OUTPUT,
			}
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
