package query

import (
	"context"
	"testing"

	gemini "google.golang.org/genai"
)

type MockGeminiClientWrapper struct {
	QueryWithContextFunc func(sessionId, userQuery, systemPrompt string) (string, error)
	QueryAgentFunc       func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error)
}

func (geminiMock *MockGeminiClientWrapper) QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error) {
	if geminiMock.QueryWithContextFunc != nil {
		return geminiMock.QueryWithContextFunc(sessionId, userQuery, systemPrompt)
	}
	return "", nil
}

func (geminiMock *MockGeminiClientWrapper) QueryAgent(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
	if geminiMock.QueryAgentFunc != nil {
		return geminiMock.QueryAgentFunc(config, history, userQuery)
	}
	return "", nil
}

func TestGeminiShouldAddSystemInstructions(t *testing.T) {
	systemInstructions := "system instruction"
	mockGeminiClientWrapper := &MockGeminiClientWrapper{
		QueryAgentFunc: func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
			return config.SystemInstruction.Parts[0].Text, nil
		},
		QueryWithContextFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
			originalGeminiClientWrapper := NewGeminiClientWrapper(nil, "model", &mockCache{}, context.Background())
			return originalGeminiClientWrapper.QueryWithContext(sessionId, userQuery, systemPrompt)
		},
	}

	queryHandler := NewQueryHandler(nil, AIClients{
		Chatgpt: &MockChatgptClient{},
		Gemini:  mockGeminiClientWrapper,
	})

	got, err := queryHandler.QueryAgent("1", Request{
		Input:    "input",
		Code:     "code",
		Language: "lang",
		Agent:    GEMINI,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != systemInstructions {
		t.Errorf("got %s, want %s", got, systemInstructions)
	}
}
