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

RULES (ABSOLUTE):
1. DO NOT USE import, require, include, or any form of library/package/module — not even standard libraries (for example: no strings, strconv, fmt, math, io, sys, etc.).
2. Only use built-in language features (variables, arrays, slices, maps, loops, conditionals, operators, primitive types).
3. Do not add package, module, or namespace declarations in the solution.
4. Do not rename the original function.
5. If the description does not make sense or does not apply, leave the code unchanged.
6. Do not write comments or explanations — output code only.
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
