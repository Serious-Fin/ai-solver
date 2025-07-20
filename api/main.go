package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"serious-fin/api/problem"
	"serious-fin/api/query"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	openai "github.com/sashabaranov/go-openai"
	gemini "google.golang.org/genai"
)

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
var chatGPTClient *openai.Client
var contextCache *query.ContextCache
var maxUserContext = 5
var problemHandler *problem.ProblemDBHandler
var queryHandler *query.QueryHandler

func main() {
	var err error

	// Initialize context cache
	contextCache, err = query.NewContextCache(maxUserContext)
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
	ctx := context.Background()
	geminiClient, err := gemini.NewClient(ctx, &gemini.ClientConfig{
		APIKey:  os.Getenv("GEMINI_KEY"),
		Backend: gemini.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Error creating gemini client: %v", err)
	}

	// Create DB handlers
	problemHandler = problem.NewProblemDBHandler(db)
	queryHandler = query.NewQueryHandler(contextCache, query.AIClients{
		Chatgpt: chatGPTClient,
		Gemini:  geminiClient,
	})

	// Initializing router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))
	router.Use(ErrorHandlerMiddleware())
	router.GET("/problems", GetProblems)
	router.GET("/problems/:id", GetProblemById)
	router.GET("/problems/:id/go", GetProblemTemplateGo)
	router.POST("/query/:sessionId", QueryAgent)
	router.POST("/validate", validateCode)

	router.Run("localhost:8080")
}

/*
TODO: Write tests for API
TODO: make authentication so not everyone could use the query endpoint to access AIs. Consider implementing a safety protocol
*/

func GetProblems(c *gin.Context) {
	problems, err := problemHandler.GetProblems()
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, problems)
}

func GetProblemById(c *gin.Context) {
	id := c.Param("id")
	problem, err := problemHandler.GetProblemById(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, problem)
}

func GetProblemTemplateGo(c *gin.Context) {
	id := c.Param("id")
	template, err := problemHandler.GetMainFuncGo(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, template)
}

func QueryAgent(c *gin.Context) {
	sessionId := c.Param("sessionId")
	var body query.Request
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	agentResponse, err := queryHandler.QueryAgent(sessionId, body)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, query.Response{
		Response: agentResponse,
	})
}
