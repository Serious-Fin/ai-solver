package query

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type ChatgptClientWrapper struct {
	Client *openai.Client
	Model  string
	Cache  *ContextCache
	Ctx    context.Context
}

func NewChatgptClientWrapper(client *openai.Client, model string, cache *ContextCache, ctx context.Context) *ChatgptClientWrapper {
	return &ChatgptClientWrapper{
		Client: client,
		Model:  model,
		Cache:  cache,
		Ctx:    ctx,
	}
}

func (wrapper *ChatgptClientWrapper) Query(sessionId string, userQuery string, systemPrompt string) (string, error) {
	previousContext := wrapper.Cache.Get(sessionId)
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

	wrapper.Cache.Add(sessionId, userQuery, output)
	return output, nil
}

func (wrapper *ChatgptClientWrapper) QueryAgent(history []any) (string, error) {
	resp, err := wrapper.Client.CreateChatCompletion(
		wrapper.Ctx,
		openai.ChatCompletionRequest{
			Model:    wrapper.Model,
			Messages: history,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
