package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

type Uri struct {
	SessionId string `uri:"sessionId" binding:"required"`
}

type QueryRequest struct {
	Input    string `form:"input"`
	Code     string `form:"code"`
	Language string `form:"language"`
	Agent    string `form:"agent"`
}

var systemPromptTemplate = `You are an expert %s programmer. The user will describe programming problems in natural 
language and also provide his current code. Respond with only %s code, no explanations, no markdown, no questions. 
You may define helper functions outside of the central function.
Do not import any external modules or packages. Code must be self-contained and runnable in a standard %s environment.`

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
	if isGoMarkdownFormat(aiOutput) {
		aiOutput, _ = strings.CutPrefix(aiOutput, "```go")
		aiOutput, _ = strings.CutSuffix(aiOutput, "```")
	}
	if isCppMarkdownFormat(aiOutput) {
		aiOutput, _ = strings.CutPrefix(aiOutput, "```cpp")
		aiOutput, _ = strings.CutSuffix(aiOutput, "```")
	}
	return aiOutput
}

func isGoMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```go") && strings.HasSuffix(str, "```")
}

func isCppMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```cpp") && strings.HasSuffix(str, "```")
}

// TODO: add gemini as another agent
