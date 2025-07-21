package query

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type ChatgptClientInterface interface {
	QueryWithContext(sessionId string, userQuery string, systemPrompt string) (string, error)
	QueryAgent(messages []openai.ChatCompletionMessage) (string, error)
}

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

func (wrapper *ChatgptClientWrapper) QueryWithContext(sessionId string, userQuery string, systemPrompt string) (string, error) {
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

	output, err := wrapper.QueryAgent(messages)
	if err != nil {
		return "", err
	}

	wrapper.Cache.Add(sessionId, userQuery, output)
	return output, nil
}

func (wrapper *ChatgptClientWrapper) QueryAgent(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := wrapper.Client.CreateChatCompletion(
		wrapper.Ctx,
		openai.ChatCompletionRequest{
			Model:    wrapper.Model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// TODO: expand errors to better explain what is wrong eg. could not query agent: %v, err
