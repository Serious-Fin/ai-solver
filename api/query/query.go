package query

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Request struct {
	Input    string `form:"input"`
	Code     string `form:"code"`
	Language string `form:"language"`
	Agent    string `form:"agent"`
}

type Response struct {
	Response string `json:"response"`
}

type AIClients struct {
	Chatgpt *openai.Client
}

type QueryHandler struct {
	Clients      AIClients
	ContextCache *ContextCache
}

func NewQueryHandler(cc *ContextCache, clients AIClients) *QueryHandler {
	return &QueryHandler{
		Clients:      clients,
		ContextCache: cc,
	}
}

var systemPrompt = `<systemPrompt>
You are an expert programmer. I need you to code solutions to programming problems. I will provide three inputs: programming language, 
current code and my own description. Description is written by me and should guide your actions. Respond only with code: no explanations, 
no markdown, no questions, no suggestions. You may define additional helper functions outside of the initial function. Do not import any 
external modules or packages.
</systemPrompt>`

var userPromptTemplate = `<description>
%s
</description>
<programmingLanguage>
%s
</programmingLanguage>
<code>
%s
</code>`

func (handler *QueryHandler) QueryAgent(sessionId string, requestBody Request) (string, error) {
	userQuery := fmt.Sprintf(userPromptTemplate, requestBody.Input, requestBody.Language, requestBody.Code)
	response, err := handler.dispatchToAgent(requestBody.Agent, sessionId, userQuery)
	if err != nil {
		return "", nil
	}
	return postProcessResponse(response), nil
}

func (handler *QueryHandler) dispatchToAgent(agent string, sessionId string, userQuery string) (string, error) {
	switch agent {
	case "chatgpt":
		return handler.queryChatGPT(sessionId, userQuery)
	default:
		return "", fmt.Errorf("agent of type %s does not exist", agent)
	}
}

func (handler *QueryHandler) queryChatGPT(sessionId string, userQuery string) (string, error) {
	previousContext := handler.ContextCache.Get(sessionId)
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

	resp, err := handler.Clients.Chatgpt.CreateChatCompletion(
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
	handler.ContextCache.Add(sessionId, userQuery, aiOutput)
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
// TODO: add tests
