package query

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type ChatgptAgentWrapperInterface interface {
	QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error)
}

type ChatgptAgentWrapper struct {
	Agent ChatgptInterface
	Cache CacheInterface
}

type ChatgptInterface interface {
	Query(messages []openai.ChatCompletionMessage) (string, error)
}

type Chatgpt struct {
	Client *openai.Client
	Model  string
	Ctx    context.Context
}

func NewChatgptClientWrapper(client *openai.Client, model string, cache *ContextCache, ctx context.Context) *ChatgptAgentWrapper {
	return &ChatgptAgentWrapper{
		Agent: &Chatgpt{
			Client: client,
			Model:  model,
			Ctx:    ctx,
		},
		Cache: cache,
	}
}

func (wrapper *ChatgptAgentWrapper) QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error) {
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

	output, err := wrapper.Agent.Query(messages)
	if err != nil {
		return "", err
	}

	wrapper.Cache.Add(sessionId, userQuery, output)
	return output, nil
}

func (wrapper *Chatgpt) Query(messages []openai.ChatCompletionMessage) (string, error) {
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
