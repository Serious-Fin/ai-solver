package query

import (
	"testing"

	"github.com/sashabaranov/go-openai"
	gemini "google.golang.org/genai"
)

type MockChatgptClient struct {
	QueryWithContextFunc func(sessionId string, userQuery string, systemPrompt string) (string, error)
	QueryAgentFunc       func(messages []openai.ChatCompletionMessage) (string, error)
}

func (gptMock *MockChatgptClient) QueryWithContext(sessionId string, userQuery string, systemPrompt string) (string, error) {
	if gptMock.QueryWithContextFunc != nil {
		return gptMock.QueryWithContextFunc(sessionId, userQuery, systemPrompt)
	}
	return "", nil
}

func (gptMock *MockChatgptClient) QueryAgent(messages []openai.ChatCompletionMessage) (string, error) {
	if gptMock.QueryAgentFunc != nil {
		return gptMock.QueryAgentFunc(messages)
	}
	return "", nil
}

type MockGeminiClient struct {
	QueryWithContextFunc func(sessionId string, userQuery string, systemPrompt string) (string, error)
	QueryAgentFunc       func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error)
}

func (geminiMock *MockGeminiClient) QueryWithContext(sessionId string, userQuery string, systemPrompt string) (string, error) {
	if geminiMock.QueryWithContextFunc != nil {
		return geminiMock.QueryWithContextFunc(sessionId, userQuery, systemPrompt)
	}
	return "", nil
}

func (geminiMock *MockGeminiClient) QueryAgent(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
	if geminiMock.QueryAgentFunc != nil {
		return geminiMock.QueryAgentFunc(config, history, userQuery)
	}
	return "", nil
}

func TestShouldInvokeChatgpt(t *testing.T) {
	want := "test code"
	mockChatgptClient := &MockChatgptClient{
		QueryWithContextFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
			return want, nil
		},
	}

	queryHandler := NewQueryHandler(nil, AIClients{
		Chatgpt: mockChatgptClient,
		Gemini:  &MockGeminiClient{},
	})

	got, err := queryHandler.QueryAgent("1", Request{
		Input:    "input",
		Code:     "code",
		Language: "lang",
		Agent:    "chatgpt",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

}
