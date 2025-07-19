package problem

import (
	"encoding/json"
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

type RowInterface interface {
	Scan(dest ...any) error
	Next() bool
	Err() error
	Close() error
}

type DBInterface interface {
	Query(query string, args ...any) (RowInterface, error)
	QueryRow(query string, args ...any) RowInterface
}

type ProblemDBHandler struct {
	DB DBInterface
}

func NewProblemDBHandler(db DBInterface) *ProblemDBHandler {
	return &ProblemDBHandler{DB: db}
}

func (handler *ProblemDBHandler) GetProblems() ([]Problem, error) {
	rows, err := handler.DB.Query("SELECT id, title FROM problems;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	problems := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		err = rows.Scan(&problem.Id, &problem.Title)
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return problems, nil
}

func (handler *ProblemDBHandler) GetProblemById(id string) (*Problem, error) {
	row := handler.DB.QueryRow("SELECT id, title, description, testCases, GoPlaceholder FROM problems WHERE id = ?;", id)

	var problem Problem
	var testCaseString string
	err := row.Scan(&problem.Id, &problem.Title, &problem.Description, &testCaseString, &problem.GoPlaceholder)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(testCaseString), &problem.TestCases)
	if err != nil {
		return nil, err
	}
	problem.TestIds = extractTestIds(problem.TestCases)
	return &problem, nil
}

func extractTestIds(testCases []TestCase) []int {
	testCaseIds := make([]int, 0)
	for _, testCase := range testCases {
		testCaseIds = append(testCaseIds, testCase.Id)
	}
	return testCaseIds
}
