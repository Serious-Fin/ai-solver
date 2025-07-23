package query

import (
	"fmt"
	"testing"
	"time"

	"github.com/sashabaranov/go-openai"
)

type MockChatgptAgentWrapper struct {
	QueryWithContextFunc func(sessionId, userQuery, systemPrompt string) (string, error)
}

func (mockWrapper *MockChatgptAgentWrapper) QueryWithContext(sessionId, userQuery, systemPrompt string) (string, error) {
	if mockWrapper.QueryWithContextFunc != nil {
		return mockWrapper.QueryWithContextFunc(sessionId, userQuery, systemPrompt)
	}
	return "", nil
}

type MockChatgpt struct {
	QueryFunc func(messages []openai.ChatCompletionMessage) (string, error)
}

func (mockChatgpt *MockChatgpt) Query(messages []openai.ChatCompletionMessage) (string, error) {
	if mockChatgpt.QueryFunc != nil {
		return mockChatgpt.QueryFunc(messages)
	}
	return "", nil
}

func TestChatgptShouldAddSystemPromptToQuery(t *testing.T) {
	chatgptAgentWrapper := &ChatgptAgentWrapper{
		Agent: &MockChatgpt{
			QueryFunc: func(messages []openai.ChatCompletionMessage) (string, error) {
				return contextToString([]Context{
					{
						Content: messages[0].Content,
						Role:    messages[0].Role,
					},
				}), nil
			},
		},
		Cache: &MockCache{},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: chatgptAgentWrapper,
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
	want := contextToString([]Context{{
		Content: systemPrompt,
		Role:    RoleSystem,
	}})
	if got != want {
		t.Errorf("got %s, want %s", got, systemPrompt)
	}
}

func TestChatgptShouldAddHistory(t *testing.T) {
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

	chatgptAgentWrapper := &ChatgptAgentWrapper{
		Agent: &MockChatgpt{
			QueryFunc: func(messages []openai.ChatCompletionMessage) (string, error) {
				context := []Context{}
				for _, ctx := range messages {
					context = append(context, Context{
						Content: ctx.Content,
						Role:    ctx.Role,
					})
				}
				return contextToString(context), nil
			},
		},
		Cache: &MockCache{
			GetFunc: func(sessionId string) []Context {
				return history
			},
		},
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: chatgptAgentWrapper,
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
	want := append([]Context{{
		Content: systemPrompt,
		Role:    RoleSystem,
	}}, history...)
	want = append(want, Context{
		Content: fmt.Sprintf(userPromptTemplate, "input", "lang", "code"),
		Role:    RoleUser,
	})
	if got != contextToString(want) {
		t.Errorf("got %s, want %s", got, contextToString(want))
	}
}

func TestChatgptHistoryShouldSavePreviousRequestsConversation(t *testing.T) {
	inputs := []string{fmt.Sprintf(userPromptTemplate, "input1", "lang1", "code1"), fmt.Sprintf(userPromptTemplate, "input2", "lang2", "code2")}
	outputs := []string{"agent response 1", "agent response 2"}
	requestResponseIndex := 0
	sessionId := "1"
	cache, _ := NewContextCache(5, time.Minute, 2*time.Minute)
	chatgptAgentWrapper := &ChatgptAgentWrapper{
		Agent: &MockChatgpt{
			QueryFunc: func(messages []openai.ChatCompletionMessage) (string, error) {
				return outputs[requestResponseIndex], nil
			},
		},
		Cache: cache,
	}

	queryHandler := NewQueryHandler(AIAgents{
		Chatgpt: chatgptAgentWrapper,
		Gemini:  &MockGeminiAgentWrapper{},
	})

	// send first message. request/response should be cached
	_, err := queryHandler.QueryAgent(sessionId, Request{
		Input:    "input1",
		Code:     "code1",
		Language: "lang1",
		Agent:    CHATGPT,
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
		Agent:    CHATGPT,
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
