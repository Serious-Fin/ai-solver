package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"serious-fin/api/problem"
	"serious-fin/api/query"
	"serious-fin/api/validator"
	"time"

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

func sendError(c *gin.Context, statusCode int, message string, err error) {
	sendDiscordMessage(err.Error())

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

			sendError(c, statusCode, message, err)
		}
	}
}

var problemHandler *problem.ProblemDBHandler
var queryHandler *query.QueryHandler
var validatorHandler *validator.ValidatorHandler

func main() {
	sessionContextSize := 5
	cacheCleanupInterval := 20 * time.Second
	sessionTimeoutInCache := 3 * time.Minute

	checkEnvVariablesOrFail()
	cache := initializeContextCacheOrFail(sessionContextSize, cacheCleanupInterval, sessionTimeoutInCache)
	database := connectToDatabaseOrFail("database.db")
	defer database.Close()
	aiHandlers := createAIAgentClientsOrFail(openai.GPT3Dot5Turbo, "gemini-2.5-flash", cache)

	problemHandler = problem.NewProblemHandler(database)
	queryHandler = query.NewQueryHandler(*aiHandlers)
	validatorHandler = validator.NewValidatorHandler(database)

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
	router.POST("/validate", ValidateCode)

	router.Run("localhost:8080")
}

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

func ValidateCode(c *gin.Context) {
	var body validator.Request
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	validatorResponse, err := validatorHandler.Validate(body)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, validatorResponse)
}

func sendDiscordMessage(message string) {
	discordToken := os.Getenv("DISCORD_TOKEN")
	channelId := os.Getenv("DISCORD_CHANNEL_ID")
	discordApiUrl := fmt.Sprintf("https://discord.com/api/channels/%s/messages", channelId)
	body := map[string]string{"content": message}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		// printing this to standard output because this function is supposed to send errors to discord normally
		log.Printf("error marshaling discord message body: %v", err)
	}
	req, err := http.NewRequest("POST", discordApiUrl, bytes.NewBuffer(bodyAsBytes))
	if err != nil {
		log.Printf("could not create new discord request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", discordToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error while sending request to discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("got response status code %d, when expected %d", resp.StatusCode, http.StatusOK)
	}
}

func initializeContextCacheOrFail(maxSize int, cleanupInterval time.Duration, sessionTimeout time.Duration) *query.ContextCache {
	contextCache, err := query.NewContextCache(maxSize, cleanupInterval, sessionTimeout)
	if err != nil {
		log.Fatalf("Error creating context cache: %v", err)
	}
	return contextCache
}

func connectToDatabaseOrFail(dbFilePath string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s", dbFilePath))
	if err != nil {
		log.Fatalf("Error while opening database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db
}

func checkEnvVariablesOrFail() {
	if envVar := os.Getenv("CHATGPT_KEY"); envVar == "" {
		log.Fatal("CHATGPT_KEY environment variable is not set")
	}
	if envVar := os.Getenv("GEMINI_KEY"); envVar == "" {
		log.Fatal("GEMINI_KEY environment variable is not set")
	}
	if envVar := os.Getenv("DISCORD_TOKEN"); envVar == "" {
		log.Fatal("DISCORD_TOKEN environment variable is not set")
	}
	if envVar := os.Getenv("DISCORD_CHANNEL_ID"); envVar == "" {
		log.Fatal("DISCORD_CHANNEL_ID environment variable is not set")
	}
}

func createAIAgentClientsOrFail(chatgptModel string, geminiModel string, cache *query.ContextCache) *query.AIAgents {
	ctx := context.Background()
	chatGPTClient := openai.NewClient(os.Getenv("CHATGPT_KEY"))
	chatgptAgent := query.NewChatgptClientWrapper(chatGPTClient, chatgptModel, cache, ctx)

	geminiClient, err := gemini.NewClient(ctx, &gemini.ClientConfig{
		APIKey:  os.Getenv("GEMINI_KEY"),
		Backend: gemini.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Error creating gemini client: %v", err)
	}
	geminiAgent := query.NewGeminiAgentWrapper(geminiClient, geminiModel, cache, ctx)
	return &query.AIAgents{
		Chatgpt: chatgptAgent,
		Gemini:  geminiAgent,
	}
}
