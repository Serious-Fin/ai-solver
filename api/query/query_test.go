package query

import (
	"testing"

	"github.com/sashabaranov/go-openai"
)

type MockChatgptClient struct {
	QueryWithContextFunc func(sessionId, userQuery, systemPrompt string) (string, error)
	QueryAgentFunc       func(messages []openai.ChatCompletionMessage) (string, error)
}

func (gptMock *MockChatgptClient) QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error) {
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

func TestShouldInvokeChatgpt(t *testing.T) {
	want := "test code"
	mockChatgptClient := &MockChatgptClient{
		QueryWithContextFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
			return want, nil
		},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: mockChatgptClient,
		Gemini:  &MockGeminiAgentWrapper{},
	})

	got, err := queryHandler.QueryAgent("1", Request{
		Input:    "input",
		Code:     "code",
		Language: "lang",
		Agent:    CHATGPT,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShouldInvokeGemini(t *testing.T) {
	want := "test code"
	mockGeminiClient := &MockGeminiAgentWrapper{
		QueryWithContextFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
			return want, nil
		},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  mockGeminiClient,
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
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShouldThrowOnUnrecognizedAgent(t *testing.T) {
	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  &MockGeminiAgentWrapper{},
	})

	_, err := queryHandler.QueryAgent("1", Request{
		Input:    "input",
		Code:     "code",
		Language: "lang",
		Agent:    "unknown",
	})
	if err == nil {
		t.Error("expected to get error, but did not get any")
	}
}

func TestShouldRemoveGoMarkdown(t *testing.T) {
	aiOutput := "```go test func ```"
	want := "test func"
	ExecuteAndExpectText(t, aiOutput, want)
}

func TestShouldRemoveCppMarkdown(t *testing.T) {
	aiOutput := "```cpp test func ```"
	want := "test func"
	ExecuteAndExpectText(t, aiOutput, want)
}

func TestShouldTrimSpace(t *testing.T) {
	aiOutput := "\n\n\t test func\t "
	want := "test func"
	ExecuteAndExpectText(t, aiOutput, want)
}

func ExecuteAndExpectText(t *testing.T, aiOutput, want string) {
	mockGeminiClient := &MockGeminiAgentWrapper{
		QueryWithContextFunc: func(sessionId, userQuery, systemPrompt string) (string, error) {
			return aiOutput, nil
		},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  mockGeminiClient,
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
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
