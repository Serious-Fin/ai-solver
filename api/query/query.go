package query

import (
	"fmt"
	"strings"

	gemini "google.golang.org/genai"
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

type AIClientWrapper interface {
	Query(sessionId string, userQuery string, systemPrompt string) (string, error)
}

type AIClients struct {
	Chatgpt AIClientWrapper
	Gemini  AIClientWrapper
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
		return handler.Clients.Chatgpt.Query(sessionId, userQuery, systemPrompt)
	case "gemini":
		return handler.Clients.Gemini.Query(sessionId, userQuery, systemPrompt)
	default:
		return "", fmt.Errorf("agent of type %s does not exist", agent)
	}
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

func getGeminiRole(role string) (gemini.Role, error) {
	switch role {
	case RoleUser:
		return gemini.RoleUser, nil
	case RoleAssistant:
		return gemini.RoleModel, nil
	default:
		return "", fmt.Errorf("unknown role type %s", role)
	}
}

// TODO: add tests
