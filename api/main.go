package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type problem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	TestCases   []testCase `json:"testCases"`
}

type testCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type APIError struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func sendAPIErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	log.Printf("API_ERROR: Status=%d, Message='%s', InternalError='%v'", statusCode, message, err)

	apiError := APIError{Message: message}
	if gin.IsDebugging() {
		apiError.Details = err.Error()
	}
	c.IndentedJSON(statusCode, apiError)
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			statusCode := http.StatusInternalServerError
			message := "An unexpected server error encountered"

			if errors.Is(err, sql.ErrNoRows) {
				statusCode = http.StatusNotFound
				message = "Resource not found"
			}

			sendAPIErrorResponse(c, statusCode, message, err)
		}
	}
}

var db *sql.DB

func main() {
	var err error

	// Connection to database
	databaseName := "database.db"
	db, err = sql.Open("sqlite3", fmt.Sprintf("./%s", databaseName))
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Successfully connected to database")
	defer db.Close()

	// Initializing router
	router := gin.Default()
	router.Use(ErrorHandlerMiddleware())
	router.GET("/problems", getProblems)
	router.GET("/problems/:id", getProblemById)

	router.Run("localhost:8080")
}

func getProblems(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, description, testCases FROM problems;")
	if err != nil {
		c.Error(err)
		return
	}
	defer rows.Close()

	problems := make([]problem, 0)
	for rows.Next() {
		var newProblem problem
		var testCasesString string
		err = rows.Scan(&newProblem.Id, &newProblem.Title, &newProblem.Description, &testCasesString)
		if err != nil {
			c.Error(err)
			return
		}
		err := json.Unmarshal([]byte(testCasesString), &newProblem.TestCases)
		if err != nil {
			c.Error(err)
			return
		}
		problems = append(problems, newProblem)
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
	row := db.QueryRow("SELECT id, title, description, testCases FROM problems WHERE id = ?;", id)

	var newProblem problem
	var testCaseString string
	err := row.Scan(&newProblem.Id, &newProblem.Title, &newProblem.Description, &testCaseString)
	if err != nil {
		c.Error(err)
		return
	}

	err = json.Unmarshal([]byte(testCaseString), &newProblem.TestCases)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, newProblem)
}
