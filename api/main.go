package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	openai "github.com/sashabaranov/go-openai"
)

type Problem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	TestCases   []TestCase `json:"testCases"`
}

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type APIError struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type QueryRequest struct {
	Input    string `form:"input"`
	Code     string `form:"code"`
	Language string `form:"language"`
	Agent    string `form:"agent"`
}

type Uri struct {
	SessionId string `uri:"sessionId" binding:"required"`
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
var chatGPTClient *openai.Client

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

	// Connect to AI agents
	chatGPTClient = openai.NewClient(os.Getenv("CHATGPT_KEY"))

	// Initializing router
	router := gin.Default()
	router.Use(ErrorHandlerMiddleware())
	router.GET("/problems", getProblems)
	router.GET("/problems/:id", getProblemById)
	router.POST("/query/:sessionId", queryAgent)

	router.Run("localhost:8080")
}

func getProblems(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, description, testCases FROM problems;")
	if err != nil {
		c.Error(err)
		return
	}
	defer rows.Close()

	problems := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		var testCasesString string
		err = rows.Scan(&problem.Id, &problem.Title, &problem.Description, &testCasesString)
		if err != nil {
			c.Error(err)
			return
		}
		err := json.Unmarshal([]byte(testCasesString), &problem.TestCases)
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
	row := db.QueryRow("SELECT id, title, description, testCases FROM problems WHERE id = ?;", id)

	var problem Problem
	var testCaseString string
	err := row.Scan(&problem.Id, &problem.Title, &problem.Description, &testCaseString)
	if err != nil {
		c.Error(err)
		return
	}

	err = json.Unmarshal([]byte(testCaseString), &problem.TestCases)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, problem)
}

func queryAgent(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.Error(err)
		return
	}

	var body QueryRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	var err error
	var resp string
	if body.Agent == "chatgpt" {
		resp, err = queryChatGPT(uri.SessionId, body.Input)
	} else {
		err = fmt.Errorf("Agent of type %s does not exist", body.Agent)
	}

	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func queryChatGPT(sessionId string, userInput string) (string, error) {
	// get user previous queries

	// make query
	resp, err := chatGPTClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert Go programmer solving LeetCode problems. Only return the code. Do not include any text, explanations, or markdown formatting outside of the code block itself. The code must be runnable. Leave main method structure as is, only modify content within or helper classes.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	// advance query

	// return response
	return resp.Choices[0].Message.Content, nil
}

/*
AI endpoint should:
1. Read body parameter which should include user code
2. Read body parameter which should include user message
3. Read body parameter of used programming language
4. Read previous context (lets say 5 last messages. Keep it in a queue in the style of system (always is)/user/assistant/user/assistent/user/assistant)
*/
