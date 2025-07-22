package query

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestGeminiShouldAddSystemPromptToQuery(t *testing.T) {
	geminiAgentWrapper := &GeminiAgentWrapper{
		Agent: &MockGemini{
			QueryFunc: func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
				return fmt.Sprintf("%s||%s", config.SystemInstruction.Parts[0].Text, config.SystemInstruction.Role), nil
			},
		},
		Cache: &MockCache{},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  geminiAgentWrapper,
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
	if got != fmt.Sprintf("%s||%s", systemPrompt, string(gemini.RoleUser)) {
		t.Errorf("got %s, want %s", got, systemPrompt)
	}
}

func TestGeminiShouldAddHistory(t *testing.T) {
	history := []Context{
		{
			Role:    RoleUser,
			Content: "Hey",
		},
		{
			Role:    RoleAssistant,
			Content: "Hello, how can I help you?",
		},
		{
			Role:    RoleUser,
			Content: "what's 2+2?",
		},
		{
			Role:    RoleAssistant,
			Content: "4",
		},
	}

	geminiAgentWrapper := &GeminiAgentWrapper{
		Agent: &MockGemini{
			QueryFunc: func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
				extractedStrings := []string{}
				for _, msg := range history {
					extractedStrings = append(extractedStrings, msg.Parts[0].Text)
					extractedStrings = append(extractedStrings, msg.Role)
				}
				return strings.Join(extractedStrings, "||"), nil
			},
		},
		Cache: &MockCache{
			GetFunc: func(sessionId string) []Context {
				return history
			},
		},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  geminiAgentWrapper,
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
	if got != contextToString(history) {
		t.Errorf("got %s, want %s", got, contextToString(history))
	}
}

func TestGeminiShouldAddUserQuery(t *testing.T) {
	wantInput := "foo bar baz"
	wantCode := "func int main"
	wantLanguage := "golang"
	want := fmt.Sprintf(userPromptTemplate, wantInput, wantLanguage, wantCode)
	geminiAgentWrapper := &GeminiAgentWrapper{
		Agent: &MockGemini{
			QueryFunc: func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
				return userQuery, nil
			},
		},
		Cache: &MockCache{},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  geminiAgentWrapper,
	})

	got, err := queryHandler.QueryAgent("1", Request{
		Input:    wantInput,
		Code:     wantCode,
		Language: wantLanguage,
		Agent:    GEMINI,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestGeminiHistoryShouldSavePreviousRequestsConversation(t *testing.T) {
	inputs := []string{fmt.Sprintf(userPromptTemplate, "input1", "lang1", "code1"), fmt.Sprintf(userPromptTemplate, "input2", "lang2", "code2")}
	outputs := []string{"agent response 1", "agent response 2"}
	requestResponseIndex := 0
	sessionId := "1"
	cache, _ := NewContextCache(5, time.Minute, 2*time.Minute)
	geminiAgentWrapper := &GeminiAgentWrapper{
		Agent: &MockGemini{
			QueryFunc: func(config *gemini.GenerateContentConfig, history []*gemini.Content, userQuery string) (string, error) {
				return outputs[requestResponseIndex], nil
			},
		},
		Cache: cache,
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: &MockChatgptClient{},
		Gemini:  geminiAgentWrapper,
	})

	// send first message. request/response should be cached
	_, err := queryHandler.QueryAgent(sessionId, Request{
		Input:    "input1",
		Code:     "code1",
		Language: "lang1",
		Agent:    GEMINI,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// send second message. request/response should be cached
	requestResponseIndex++
	_, err = queryHandler.QueryAgent(sessionId, Request{
		Input:    "input2",
		Code:     "code2",
		Language: "lang2",
		Agent:    GEMINI,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// get the history and compare it with expected (two request and two responses cached)
	history := cache.Get(sessionId)
	if len(history) != 4 {
		t.Errorf("expected to have 4 messages in cache but found: %d", len(history))
	}

	want := contextToString([]Context{
		{
			Content: inputs[0],
			Role:    RoleUser,
		},
		{
			Content: outputs[0],
			Role:    RoleAssistant,
		},
		{
			Content: inputs[1],
			Role:    RoleUser,
		},
		{
			Content: outputs[1],
			Role:    RoleAssistant,
		},
	})
	got := contextToString(history)
	if want != got {
		t.Errorf("got %s, want %s", got, want)
	}
}

func contextToString(context []Context) string {
	extractedStrings := []string{}
	for _, msg := range context {
		extractedStrings = append(extractedStrings, msg.Content)
		role, _ := getGeminiRole(msg.Role)
		extractedStrings = append(extractedStrings, string(role))
	}
	return strings.Join(extractedStrings, "||")
}
