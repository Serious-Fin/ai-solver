package query

import (
	"sync"
	"testing"
	"time"
)

type mockCache struct {
	AddFunc                 func(sessionId, userInput, aiOutput string)
	GetFunc                 func(sessionId string) []Context
	StartCleanupRoutineFunc func()
	StopCleanupRoutineFunc  func()
}

func (mockCache *mockCache) Add(sessionId, userInput, aiOutput string) {
	if mockCache.AddFunc != nil {
		mockCache.AddFunc(sessionId, userInput, aiOutput)
	}
}

func (mockCache *mockCache) Get(sessionId string) []Context {
	if mockCache.GetFunc != nil {
		return mockCache.GetFunc(sessionId)
	}
	return nil
}

func (mockCache *mockCache) StartCleanupRoutine() {
	if mockCache.StartCleanupRoutineFunc != nil {
		mockCache.StartCleanupRoutineFunc()
	}
}

func (mockCache *mockCache) StopCleanupRoutine() {
	if mockCache.StopCleanupRoutineFunc != nil {
		mockCache.StopCleanupRoutineFunc()
	}
}

func newMockTime() *MockTime {
	return &MockTime{currTime: time.Now()}
}

type MockTime struct {
	mu       sync.Mutex
	currTime time.Time
}

func (mockTime *MockTime) Now() time.Time {
	mockTime.mu.Lock()
	defer mockTime.mu.Unlock()
	return mockTime.currTime
}

func (mockTime *MockTime) Advance(duration time.Duration) {
	mockTime.mu.Lock()
	defer mockTime.mu.Unlock()
	mockTime.currTime = mockTime.currTime.Add(duration)
}

func TestConstructorErrorMaxSize(t *testing.T) {
	_, err := NewContextCache(-4, time.Second, time.Minute)
	if err == nil {
		t.Error("Expected constructor with negative maxSize to throw")
	}
}

func TestConstructorErrorCleanupInterval(t *testing.T) {
	_, err := NewContextCache(3, -time.Second, time.Minute)
	if err == nil {
		t.Error("Expected constructor with negative cleanupInterval to throw")
	}
}

func TestConstructorErrorSessionTimeout(t *testing.T) {
	_, err := NewContextCache(3, time.Second, -time.Minute)
	if err == nil {
		t.Error("Expected constructor with negative sessionTimeout to throw")
	}
}

func TestCacheAddAndGetSameSession(t *testing.T) {
	expect := []Context{
		{
			Content: "input1",
			Role:    RoleUser,
		},
		{
			Content: "output1",
			Role:    RoleAssistant,
		},
		{
			Content: "input2",
			Role:    RoleUser,
		},
		{
			Content: "output2",
			Role:    RoleAssistant,
		},
	}
	cache, err := NewContextCache(5, time.Second, time.Minute)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}

	cache.Add("1", expect[0].Content, expect[1].Content)
	cache.Add("1", expect[2].Content, expect[3].Content)
	history := cache.Get("1")

	got := contextToString(history)
	want := contextToString(expect)
	if got != want {
		t.Errorf("got %v, want: %v", got, want)
	}
}

func TestCacheAddAndGetDifferentSessions(t *testing.T) {
	cache, err := NewContextCache(5, time.Second, time.Minute)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}

	cache.Add("1", "input", "output")
	history := cache.Get("2")

	got := contextToString(history)
	if got != "" {
		t.Errorf("expected empty string but got %v", got)
	}
}

func TestCacheExceedsLimit(t *testing.T) {
	expect := []Context{
		{
			Content: "input0",
			Role:    RoleUser,
		},
		{
			Content: "output0",
			Role:    RoleAssistant,
		},
		{
			Content: "input1",
			Role:    RoleUser,
		},
		{
			Content: "output1",
			Role:    RoleAssistant,
		},
		{
			Content: "input2",
			Role:    RoleUser,
		},
		{
			Content: "output2",
			Role:    RoleAssistant,
		},
	}
	cache, err := NewContextCache(2, time.Second, time.Minute)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}

	cache.Add("1", expect[0].Content, expect[1].Content)
	cache.Add("1", expect[2].Content, expect[3].Content)
	history := cache.Get("1")
	got := contextToString(history)
	want := contextToString(expect[:4])
	if got != want {
		t.Errorf("got %v, want: %v", got, want)
	}

	cache.Add("1", expect[4].Content, expect[5].Content)
	history = cache.Get("1")
	got = contextToString(history)
	want = contextToString(expect[2:])
	if got != want {
		t.Errorf("got %v, want: %v", got, want)
	}
}

func TestCacheCleansUpStagnantSessions(t *testing.T) {
	timer := newMockTime()
	cleanupInterval := 10 * time.Millisecond
	sessionTimeout := 30 * time.Millisecond
	cache, err := NewContextCacheWithTimeFunc(3, cleanupInterval, sessionTimeout, timer.Now)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}
	cache.StartCleanupRoutine()

	cache.Add("1", "input1", "output1")
	timer.Advance(10 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got := cache.Get("1")
	if len(got) != 2 {
		t.Errorf("got %d, want %d context items", len(got), 2)
	}

	timer.Advance(30 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got = cache.Get("1")
	if len(got) != 0 {
		t.Errorf("got %d, want %d context items", len(got), 0)
	}
}

func TestCacheAddResetsTimeout(t *testing.T) {
	timer := newMockTime()
	cleanupInterval := 10 * time.Millisecond
	sessionTimeout := 30 * time.Millisecond
	cache, err := NewContextCacheWithTimeFunc(3, cleanupInterval, sessionTimeout, timer.Now)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}
	cache.StartCleanupRoutine()

	cache.Add("1", "input1", "output1")
	timer.Advance(20 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)

	cache.Add("1", "input2", "output2")
	timer.Advance(20 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)

	cache.Add("1", "input3", "output3")
	timer.Advance(20 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got := cache.Get("1")
	if len(got) != 6 {
		t.Errorf("got %d, want %d context items", len(got), 6)
	}
}

func TestCacheGetResetsTimeout(t *testing.T) {
	timer := newMockTime()
	cleanupInterval := 10 * time.Millisecond
	sessionTimeout := 30 * time.Millisecond
	cache, err := NewContextCacheWithTimeFunc(3, cleanupInterval, sessionTimeout, timer.Now)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}
	cache.StartCleanupRoutine()

	cache.Add("1", "input1", "output1")

	timer.Advance(20 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got := cache.Get("1")
	if len(got) != 2 {
		t.Errorf("got %d, want %d context items", len(got), 2)
	}

	timer.Advance(20 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got = cache.Get("1")
	if len(got) != 2 {
		t.Errorf("got %d, want %d context items", len(got), 2)
	}

	timer.Advance(20 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got = cache.Get("1")
	if len(got) != 2 {
		t.Errorf("got %d, want %d context items", len(got), 2)
	}
}

func TestCacheStopCleanupRoutine(t *testing.T) {
	timer := newMockTime()
	cleanupInterval := 10 * time.Millisecond
	sessionTimeout := 30 * time.Millisecond
	cache, err := NewContextCacheWithTimeFunc(3, cleanupInterval, sessionTimeout, timer.Now)
	if err != nil {
		t.Errorf("Unexpected error while creating cache: %v", err)
	}
	cache.StartCleanupRoutine()

	cache.Add("1", "input1", "output1")
	timer.Advance(30 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got := cache.Get("1")
	if len(got) != 0 {
		t.Errorf("got %d, want %d context items", len(got), 0)
	}

	cache.StopCleanupRoutine()
	timer.Advance(30 * time.Millisecond)

	cache.Add("1", "input2", "output2")
	timer.Advance(30 * time.Millisecond)
	time.Sleep(cleanupInterval * 2)
	got = cache.Get("1")
	if len(got) != 2 {
		t.Errorf("got %d, want %d context items", len(got), 2)
	}
}
