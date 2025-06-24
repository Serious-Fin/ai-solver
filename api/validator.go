package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidateRequest struct {
	ProblemId string `form:"problemId"`
	Code      string `form:"code"`
	Language  string `form:"language"`
}

type ValidateResponse struct {
	FailedTests    []int `json:"failed"`
	SucceededTests []int `json:"succeeded"`
}

type TestParams struct {
	Template string
	Helpers  string
	Cases    []TestCase
}

var fileStartTemplate = `package main
import "testing"
`

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
	// validate
	// run the following:
	//  docker run -it --rm -v ./test_run_2531706648:/app --network none go-testing-image:latest /bin/sh -c "go mod init test_proj && go test -v"
	// remove temp folder
	c.IndentedJSON(http.StatusOK, ValidateResponse{})
}

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
		newTestCode = strings.Replace(newTestCode, "{{ID}}", strconv.Itoa(testCase.Id), 1)
		newTestCode = strings.Replace(newTestCode, "{{OUTPUT}}", testCase.ExpectedOutput, 1)
		fmt.Printf("%v\n", testCase.Inputs)
		for inputIndex, input := range testCase.Inputs {
			toChange := fmt.Sprintf("{{INPUT%d}}", inputIndex)
			fmt.Printf("Changing %s with %s", toChange, input)
			newTestCode = strings.Replace(newTestCode, fmt.Sprintf("{{INPUT%d}}", inputIndex), input, 1)
		}
		_, err = file.WriteString(fmt.Sprintf("%s\n", newTestCode))
		check(err)
	}

	_, err = file.WriteString(helperFuncs)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
