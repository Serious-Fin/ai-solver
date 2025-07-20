package query

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
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

type AIClients struct {
	Chatgpt *openai.Client
	Gemini  *gemini.Client
}

type QueryHandler struct {
	Clients      AIClients
	Context      context.Context
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
	case "gemini":
		return handler.queryGemini(sessionId, userQuery)
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

func (handler *QueryHandler) queryGemini(sessionId string, userQuery string) (string, error) {
	config := &gemini.GenerateContentConfig{
		SystemInstruction: gemini.NewContentFromText(systemPrompt, gemini.RoleUser),
	}

	previousContext := handler.ContextCache.Get(sessionId)
	history := make([]*gemini.Content, 0)
	for _, context := range previousContext {
		role, err := getGeminiRole(context.Role)
		if err != nil {
			return "", err
		}
		history = append(history, gemini.NewContentFromText(context.Content, role))
	}

	chat, err := handler.Clients.Gemini.Chats.Create(handler.Context, "gemini-2.5-flash", config, history)
	if err != nil {
		return "", fmt.Errorf("failed to initialize new gemini chat session: %v", err)
	}
	res, err := chat.SendMessage(handler.Context, gemini.Part{Text: userQuery})
	if err != nil {
		return "", fmt.Errorf("failed to send new message to gemini: %v", err)
	}

	if len(res.Candidates) == 0 {
		return "", fmt.Errorf("no response was received from gemini query")
	}
	return res.Candidates[0].Content.Parts[0].Text, nil
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
