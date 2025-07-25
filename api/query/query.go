package query

import (
	"fmt"
	"strings"
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

type AIAgents struct {
	Chatgpt ChatgptAgentWrapperInterface
	Gemini  GeminiAgentWrapperInterface
}

type QueryHandler struct {
	Agents AIAgents
}

func NewQueryHandler(agents AIAgents) *QueryHandler {
	return &QueryHandler{
		Agents: agents,
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

const (
	GEMINI  = "gemini"
	CHATGPT = "chatgpt"
)

func (handler *QueryHandler) QueryAgent(sessionId string, requestBody Request) (string, error) {
	userQuery := fmt.Sprintf(userPromptTemplate, requestBody.Input, requestBody.Language, requestBody.Code)
	response, err := handler.dispatchToAgent(requestBody.Agent, sessionId, userQuery)
	if err != nil {
		return "", fmt.Errorf("error querying agent: %v", err)
	}
	return postProcessResponse(response), nil
}

func (handler *QueryHandler) dispatchToAgent(agent, sessionId, userQuery string) (string, error) {
	switch agent {
	case CHATGPT:
		return handler.Agents.Chatgpt.QueryWithContext(sessionId, userQuery, systemPrompt)
	case GEMINI:
		return handler.Agents.Gemini.QueryWithContext(sessionId, userQuery, systemPrompt)
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
	aiOutput = strings.TrimSpace(aiOutput)
	return aiOutput
}

func isGoMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```go") && strings.HasSuffix(str, "```")
}

func isCppMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```cpp") && strings.HasSuffix(str, "```")
}
