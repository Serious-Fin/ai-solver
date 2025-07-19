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
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	openai "github.com/sashabaranov/go-openai"
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
var contextCache *ContextCache
var maxUserContext = 5
var systemPromptTemplate = `You are an expert %s programmer. The user will describe programming problems in natural 
language and also provide his current code. Respond with only %s code, no explanations, no markdown, no questions. 
Do not modify the function named run. It is the entry point for evaluation. You may define helper functions outside 
run. Do not import any external modules or packages. Code must be self-contained and runnable in a standard %s environment.`

func main() {
	var err error

	// Initialize context cache
	contextCache, err = NewContextCache(maxUserContext)
	if err != nil {
		log.Fatalf("Error creating context cache: %v", err)
	}

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
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))
	router.Use(ErrorHandlerMiddleware())
	router.GET("/problems", getProblems)
	router.GET("/problems/:id", getProblemById)
	router.POST("/query/:sessionId", queryAgent)
	router.POST("/validate", validateCode)

	router.Run("localhost:8080")
}

func getProblems(c *gin.Context) {
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
	systemPrompt := fmt.Sprintf(systemPromptTemplate, body.Language, body.Language, body.Language)
	userQuery := fmt.Sprintf("User input: %s\nProgramming language used: %s\nCurrent user code: %s", body.Input, body.Language, body.Code)
	if body.Agent == "chatgpt" {
		resp, err = queryChatGPT(uri.SessionId, systemPrompt, userQuery)
	} else {
		err = fmt.Errorf("agent of type %s does not exist", body.Agent)
	}

	if err != nil {
		c.Error(err)
		return
	}

	resp = postProcessResponse(resp)

	c.IndentedJSON(http.StatusOK, struct {
		Response string `json:"response"`
	}{Response: resp})
}

func queryChatGPT(sessionId string, systemPrompt string, userQuery string) (string, error) {
	previousContext := contextCache.Get(sessionId)
	messages := make([]openai.ChatCompletionMessage, 0)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    RoleSystem,
		Content: systemPrompt,
	})
	for _, context := range previousContext {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    context.Role,
			Content: context.Content,
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    RoleUser,
		Content: userQuery,
	})

	resp, err := chatGPTClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}

	aiOutput := resp.Choices[0].Message.Content
	contextCache.Add(sessionId, userQuery, aiOutput)
	return aiOutput, nil
}

func postProcessResponse(aiOutput string) string {
	if isMarkdownFormat(aiOutput) {
		aiOutput, _ = strings.CutPrefix(aiOutput, "```go")
		aiOutput, _ = strings.CutSuffix(aiOutput, "```")
	}
	return aiOutput
}

func isMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```go") && strings.HasSuffix(str, "```")
}

/*
TODO: Write tests for API
TODO: make authentication so not everyone could use the query endpoint to access AIs
TODO: remove ```cpp
*/
