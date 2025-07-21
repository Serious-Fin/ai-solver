package query

import (
	"testing"
)

type MockAIClient struct {
	QueryFunc func(sessionId string, userQuery string, systemPrompt string) (string, error)
}

func (gptMock *MockAIClient) Query(sessionId string, userQuery string, systemPrompt string) (string, error) {
	if gptMock.QueryFunc != nil {
		return gptMock.QueryFunc(sessionId, userQuery, systemPrompt)
	}
	return "", nil
}

func TestQueryChatgpt(t *testing.T) {
	mockChatgptClient := &MockAIClient{
		QueryFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
			return "test code", nil
		},
	}

	queryHandler := NewQueryHandler(nil, AIClients{
		Chatgpt: mockChatgptClient,
		Gemini:  &MockAIClient{},
	})

}
