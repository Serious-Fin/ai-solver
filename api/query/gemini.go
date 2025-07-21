package query

import (
	"context"
	"fmt"

	gemini "google.golang.org/genai"
)

type GeminiCLientInterface interface {
	QueryWithContext(sessionId string, userQuery string, systemPrompt string) (string, error)
	QueryAgent(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error)
}

type GeminiClientWrapper struct {
	Client *gemini.Client
	Model  string
	Cache  *ContextCache
	Ctx    context.Context
}

func NewGeminiClientWrapper(client *gemini.Client, model string, cache *ContextCache, ctx context.Context) *GeminiClientWrapper {
	return &GeminiClientWrapper{
		Client: client,
		Model:  model,
		Cache:  cache,
		Ctx:    ctx,
	}
}

func (wrapper *GeminiClientWrapper) QueryWithContext(sessionId string, userQuery string, systemPrompt string) (string, error) {
	config := &gemini.GenerateContentConfig{
		SystemInstruction: gemini.NewContentFromText(systemPrompt, gemini.RoleUser),
	}

	previousContext := wrapper.Cache.Get(sessionId)
	history := make([]*gemini.Content, 0)
	for _, context := range previousContext {
		role, err := getGeminiRole(context.Role)
		if err != nil {
			return "", err
		}
		history = append(history, gemini.NewContentFromText(context.Content, role))
	}

	output, err := wrapper.QueryAgent(config, history, userQuery)
	if err != nil {
		return "", nil
	}
	wrapper.Cache.Add(sessionId, userQuery, output)
	return output, nil
}

func (wrapper *GeminiClientWrapper) QueryAgent(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
	chat, err := wrapper.Client.Chats.Create(wrapper.Ctx, wrapper.Model, config, history)
	if err != nil {
		return "", fmt.Errorf("failed to initialize new gemini chat session: %v", err)
	}
	res, err := chat.SendMessage(wrapper.Ctx, gemini.Part{Text: userQuery})
	if err != nil {
		return "", fmt.Errorf("failed to send new message to gemini: %v", err)
	}

	if len(res.Candidates) == 0 {
		return "", fmt.Errorf("no response was received from gemini query")
	}
	return res.Candidates[0].Content.Parts[0].Text, nil
}
