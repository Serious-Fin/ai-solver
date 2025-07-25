package query

import (
	"context"
	"fmt"

	gemini "google.golang.org/genai"
)

type GeminiAgentWrapperInterface interface {
	QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error)
}

type GeminiAgentWrapper struct {
	Agent GeminiInterface
	Cache CacheInterface
}

type GeminiInterface interface {
	Query(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error)
}

type Gemini struct {
	Client *gemini.Client
	Model  string
	Ctx    context.Context
}

func NewGeminiAgentWrapper(client *gemini.Client, model string, cache CacheInterface, ctx context.Context) *GeminiAgentWrapper {
	return &GeminiAgentWrapper{
		Agent: &Gemini{
			Client: client,
			Model:  model,
			Ctx:    ctx,
		},
		Cache: cache,
	}
}

func (wrapper *GeminiAgentWrapper) QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error) {
	config := &gemini.GenerateContentConfig{
		SystemInstruction: gemini.NewContentFromText(systemPrompt, gemini.RoleUser),
	}

	previousContext := wrapper.Cache.Get(sessionId)
	history := make([]*gemini.Content, 0)
	for _, context := range previousContext {
		role, err := getGeminiRole(context.Role)
		if err != nil {
			return "", fmt.Errorf("error building previous context for gemini request: %v", err)
		}
		history = append(history, gemini.NewContentFromText(context.Content, role))
	}

	output, err := wrapper.Agent.Query(config, history, userQuery)
	if err != nil {
		return "", fmt.Errorf("could not query gemini agent: %v", err)
	}
	wrapper.Cache.Add(sessionId, userQuery, output)
	return output, nil
}

func (agent *Gemini) Query(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
	chat, err := agent.Client.Chats.Create(agent.Ctx, agent.Model, config, history)
	if err != nil {
		return "", fmt.Errorf("failed to initialize new gemini chat session: %v", err)
	}
	res, err := chat.SendMessage(agent.Ctx, gemini.Part{Text: userQuery})
	if err != nil {
		return "", fmt.Errorf("failed to send new message to gemini: %v", err)
	}

	if len(res.Candidates) == 0 {
		return "", fmt.Errorf("no response was received from gemini query")
	}
	return res.Candidates[0].Content.Parts[0].Text, nil
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
