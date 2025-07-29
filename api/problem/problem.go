package problem

import (
	"encoding/json"
	"fmt"
	"serious-fin/api/common"
)

type Problem struct {
	Id            int               `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description,omitempty"`
	TestCases     []common.TestCase `json:"testCases,omitempty"`
	GoPlaceholder string            `json:"goPlaceholder,omitempty"`
}

type ProblemDBHandler struct {
	DB common.DBInterface
}

func NewProblemHandler(db common.DBInterface) *ProblemDBHandler {
	return &ProblemDBHandler{DB: db}
}

func (handler *ProblemDBHandler) GetProblems() ([]Problem, error) {
	rows, err := handler.DB.Query("SELECT id, title FROM problems")
	if err != nil {
		return nil, fmt.Errorf("could not query problems data from db: %v", err)
	}
	defer rows.Close()

	problems := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		err = rows.Scan(&problem.Id, &problem.Title)
		if err != nil {
			return nil, fmt.Errorf("could not scan problems db output: %v", err)
		}
		problems = append(problems, problem)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error reading db output: %v", err)
	}
	return problems, nil
}

func (handler *ProblemDBHandler) GetProblemById(id string) (*Problem, error) {
	row := handler.DB.QueryRow("SELECT id, title, description, testCases FROM problems WHERE id = ?", id)

	var problem Problem
	var testCaseString string
	err := row.Scan(&problem.Id, &problem.Title, &problem.Description, &testCaseString)
	if err != nil {
		return nil, fmt.Errorf("could not scan single problem db output (problem id %s): %v", id, err)
	}

	err = json.Unmarshal([]byte(testCaseString), &problem.TestCases)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal test cases into object (problem id %s): %v", id, err)
	}
	return &problem, nil
}

func (handler *ProblemDBHandler) GetMainFuncGo(problemId string) (string, error) {
	row := handler.DB.QueryRow("SELECT mainFunction FROM goTemplates WHERE problemFk = ?", problemId)

	var mainFunction string
	err := row.Scan(&mainFunction)
	if err != nil {
		return "", fmt.Errorf("could not scan problem template db output (problem id %s): %v", problemId, err)
	}

	return mainFunction, nil
}
