package query

import (
	"testing"

	gemini "google.golang.org/genai"
)

type MockGeminiAgentWrapper struct {
	QueryWithContextFunc func(sessionId, userQuery, systemPrompt string) (string, error)
}

func (mockWrapper *MockGeminiAgentWrapper) QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error) {
	if mockWrapper.QueryWithContextFunc != nil {
		return mockWrapper.QueryWithContextFunc(sessionId, userQuery, systemPrompt)
	}
	return "", nil
}

type MockGemini struct {
	QueryFunc func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error)
}

func (mockGemini *MockGemini) Query(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
	if mockGemini.QueryFunc != nil {
		return mockGemini.QueryFunc(config, history, userQuery)
	}
	return "", nil
}

func TestGeminiShouldAddSystemInstructions(t *testing.T) {
	/*
		systemInstructions := "system instruction"
		mockGeminiClientWrapper := &MockGeminiClientWrapper{
			QueryAgentFunc: func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
				return config.SystemInstruction.Parts[0].Text, nil
			},
			QueryWithContextFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
				originalGeminiClientWrapper := NewGeminiAgentWrapper(nil, "model", &mockCache{}, context.Background())
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
	*/
}
