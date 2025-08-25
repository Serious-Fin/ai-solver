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
current code, and my own description. The description is written by me and should guide your actions. Respond only with code: no explanations, 
no markdown, no questions, no suggestions. You may define additional helper functions outside of the initial function.

RULES (STRICT):
1. Standard libraries ARE allowed, but if you use them, you MUST import them explicitly at the top of the file. 
   (For example: in Go, always use import "strings" before using strings.Builder.)
2. Do not use non-standard or third-party packages — only built-in standard libraries.
3. Never omit imports if a library is used.
4. Never change the name of the original function (this function is always at the top of the file).
5. Never declare a package/module/namespace beyond the required imports (e.g., no package main in Go).
6. If user input is nonsense and you cannot change the code meaningfully, do nothing at all.
7. Do not communicate with the user — no comments, no explanations, only code.
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
		return "", fmt.Errorf("error querying agent: %w", err)
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
	if isCodeXML(aiOutput) {
		aiOutput, _ = strings.CutPrefix(aiOutput, "<code>")
		aiOutput, _ = strings.CutSuffix(aiOutput, "</code>")
	}
	aiOutput = strings.TrimSpace(aiOutput)
	return aiOutput
}

func isCodeXML(str string) bool {
	return strings.HasPrefix(str, "<code>") && strings.HasSuffix(str, "</code>")
}

func isGoMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```go") && strings.HasSuffix(str, "```")
}

func isCppMarkdownFormat(str string) bool {
	return strings.HasPrefix(str, "```cpp") && strings.HasSuffix(str, "```")
}
