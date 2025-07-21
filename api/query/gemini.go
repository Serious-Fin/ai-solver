package query

import (
	"context"
	"fmt"

	gemini "google.golang.org/genai"
)

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

func (wrapper *GeminiClientWrapper) Query(sessionId string, userQuery string, systemPrompt string) (string, error) {
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
	output := res.Candidates[0].Content.Parts[0].Text
	wrapper.Cache.Add(sessionId, userQuery, output)
	return output, nil
}
