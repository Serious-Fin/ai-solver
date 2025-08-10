package problem

import (
	"encoding/json"
	"fmt"
	"serious-fin/api/common"
)

type Problem struct {
	Id            int               `json:"id"`
	Title         string            `json:"title"`
	Difficulty    int               `json:"difficulty"`
	IsCompleted   bool              `json:"isCompleted"`
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

func (handler *ProblemDBHandler) GetProblems(userId string) ([]Problem, error) {
	query := `
	SELECT 
		id, 
		title, 
		difficulty, 
		CASE WHEN ucp.problemId IS NULL 
			THEN false 
			ELSE true 
		END AS isCompleted 
	FROM problems 
	LEFT JOIN (
		SELECT 
			problemId 
		FROM userCompletedProblems 
		WHERE userId = ?
	) AS ucp 
	ON problems.id = ucp.problemId`

	rows, err := handler.DB.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("could not query problems data from db: %w", err)
	}
	defer rows.Close()

	problems := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		err = rows.Scan(&problem.Id, &problem.Title, &problem.Difficulty, &problem.IsCompleted)
		if err != nil {
			return nil, fmt.Errorf("could not scan problems db output: %w", err)
		}
		problems = append(problems, problem)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error reading db output: %w", err)
	}
	return problems, nil
}

func (handler *ProblemDBHandler) GetProblemById(userId, problemId string) (*Problem, error) {
	query := `
	SELECT 
		id, 
		title, 
		difficulty, 
		description, 
		testCases,
		CASE WHEN ucp.problemId IS NULL 
			THEN false 
			ELSE true 
		END AS isCompleted 
	FROM problems
	LEFT JOIN (
		SELECT 
			problemId 
		FROM userCompletedProblems 
		WHERE userId = ?
		AND problemId = ?
	) AS ucp 
	ON problems.id = ucp.problemId
	WHERE problems.id = ?`

	row := handler.DB.QueryRow(query, userId, problemId, problemId)
	var problem Problem
	var testCaseString string
	err := row.Scan(&problem.Id, &problem.Title, &problem.Difficulty, &problem.Description, &testCaseString, &problem.IsCompleted)
	if err != nil {
		return nil, fmt.Errorf("could not scan single problem db output (problem id %s): %w", problemId, err)
	}

	err = json.Unmarshal([]byte(testCaseString), &problem.TestCases)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal test cases into object (problem id %s): %w", problemId, err)
	}
	return &problem, nil
}

func (handler *ProblemDBHandler) GetMainFuncGo(problemId string) (string, error) {
	row := handler.DB.QueryRow("SELECT mainFunction FROM goTemplates WHERE problemFk = ?", problemId)

	var mainFunction string
	err := row.Scan(&mainFunction)
	if err != nil {
		return "", fmt.Errorf("could not scan problem template db output (problem id %s): %w", problemId, err)
	}

	return mainFunction, nil
}
