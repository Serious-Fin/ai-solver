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
	"serious-fin/api/user"
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
	sendToDiscord(err.Error())

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
var userHandler *user.UserDBHandler

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
	userHandler = user.NewUserHandler(database)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	router.Use(ErrorHandlerMiddleware())
	router.GET("/problems", GetProblems)
	router.GET("/problems/:id", GetProblemById)
	router.POST("/problems/:id", CompleteProblem)
	router.GET("/problems/:id/go", GetProblemTemplateGo)
	router.POST("/query/:sessionId", QueryAgent)
	router.POST("/validate", ValidateCode)
	router.GET("/user/:userId", GetUser)
	router.POST("/user", CreateUser)
	router.POST("/session", StartSession)
	router.GET("/session/:sessionId", GetSession)

	router.Run("localhost:8080")
}

func GetProblems(c *gin.Context) {
	userId := c.DefaultQuery("user", "0")
	problems, err := problemHandler.GetProblems(userId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, problems)
}

func GetProblemById(c *gin.Context) {
	problemId := c.Param("id")
	userId := c.DefaultQuery("user", "0")
	problem, err := problemHandler.GetProblemById(userId, problemId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, problem)
}

func CompleteProblem(c *gin.Context) {
	problemId := c.Param("id")
	var body problem.CompleteProblemRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	err := problemHandler.CompleteProblem(problemId, body.UserId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
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

func GetSession(c *gin.Context) {
	sessionId := c.Param("sessionId")
	foundUser, err := userHandler.GetUserFromSession(sessionId)
	if err != nil {
		c.Error(err)
		return
	}

	if foundUser != nil {
		c.IndentedJSON(http.StatusOK, user.SessionInfoResponse{
			User: *foundUser,
		})
		return
	}

	c.IndentedJSON(http.StatusNotFound, nil)
}

func GetUser(c *gin.Context) {
	userId := c.Param("userId")
	existingUser, err := userHandler.GetUser(userId)
	if err != nil {
		c.Error(err)
		return
	}

	if existingUser == nil {
		c.IndentedJSON(http.StatusNotFound, nil)
		return
	}

	c.IndentedJSON(http.StatusOK, user.UserResponse{
		User: *existingUser,
	})
}

func CreateUser(c *gin.Context) {
	var body user.User
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}
	newUser, err := userHandler.CreateUser(body)
	if err != nil {
		c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, user.UserResponse{
		User: *newUser,
	})
}

func StartSession(c *gin.Context) {
	var body user.SessionRequest
	if err := c.ShouldBind(&body); err != nil {
		c.Error(err)
		return
	}

	existingSession, err := userHandler.GetSession(body.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	if existingSession != nil && !user.IsSessionExpired(existingSession) {
		extendedSession, err := userHandler.UpdateSession(existingSession.Id)
		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusOK, user.SessionResponse{
			SessionId: extendedSession.Id,
		})
		return
	}

	err = userHandler.CleanupExpiredSessions(body.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	newSession, err := userHandler.CreateSession(body.UserId)
	if err != nil {
		c.Error(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, user.SessionResponse{
		SessionId: newSession.Id,
	})
}

func sendToDiscord(message string) {
	runes := []rune(message)
	sentChars := 0
	for sentChars < len(runes) {
		end := min(sentChars+2000, len(runes))
		sendDiscordMessage(string(runes[sentChars:end]))
		sentChars = end
	}
}

func sendDiscordMessage(message string) {
	discordToken := os.Getenv("DISCORD_TOKEN")
	channelId := os.Getenv("DISCORD_CHANNEL_ID")
	discordApiUrl := fmt.Sprintf("https://discord.com/api/channels/%s/messages", channelId)
	body := map[string]string{"content": fmt.Sprintf("backend-msg: %s", message)}
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		// printing this to standard output because this function is supposed to send errors to discord normally
		fmt.Printf("error marshaling discord message body: %v", err)
	}
	req, err := http.NewRequest("POST", discordApiUrl, bytes.NewBuffer(bodyAsBytes))
	if err != nil {
		fmt.Printf("could not create new discord request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", discordToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error while sending request to discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("got response status code %d, when expected %d", resp.StatusCode, http.StatusOK)
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
