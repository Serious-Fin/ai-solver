package problems

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Problem struct {
	Id            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	TestCases     []TestCase `json:"testCases,omitempty"`
	GoPlaceholder string     `json:"goPlaceholder,omitempty"`
	TestIds       []int      `json:"testCaseIds,omitempty"`
}

type TestCase struct {
	Id             int      `json:"id"`
	Inputs         []string `json:"inputs"`
	ExpectedOutput string   `json:"output"`
}

func GetProblems(c *gin.Context) {
	rows, err := db.Query("SELECT id, title FROM problems;")
	if err != nil {
		c.Error(err)
		return
	}
	defer rows.Close()

	problems := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		err = rows.Scan(&problem.Id, &problem.Title)
		if err != nil {
			c.Error(err)
			return
		}
		problems = append(problems, problem)
	}
	err = rows.Err()
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, problems)
}

func getProblemById(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("SELECT id, title, description, testCases, GoPlaceholder FROM problems WHERE id = ?;", id)

	var problem Problem
	var testCaseString string
	err := row.Scan(&problem.Id, &problem.Title, &problem.Description, &testCaseString, &problem.GoPlaceholder)
	if err != nil {
		c.Error(err)
		return
	}

	err = json.Unmarshal([]byte(testCaseString), &problem.TestCases)
	if err != nil {
		c.Error(err)
		return
	}
	problem.TestIds = extractTestIds(problem.TestCases)
	c.IndentedJSON(http.StatusOK, problem)
}

func extractTestIds(testCases []TestCase) []int {
	testCaseIds := make([]int, 0)
	for _, testCase := range testCases {
		testCaseIds = append(testCaseIds, testCase.Id)
	}
	return testCaseIds
}
