package query

type MockCache struct {
	AddFunc                 func(sessionId, userInput, aiOutput string)
	GetFunc                 func(sessionId string) []Context
	StartCleanupRoutineFunc func()
	StopCleanupRoutineFunc  func()
}

func (mockCache *MockCache) Add(sessionId, userInput, aiOutput string) {
	if mockCache.AddFunc != nil {
		mockCache.AddFunc(sessionId, userInput, aiOutput)
	}
}

func (mockCache *MockCache) Get(sessionId string) []Context {
	if mockCache.GetFunc != nil {
		return mockCache.GetFunc(sessionId)
	}
	return nil
}

func (mockCache *MockCache) StartCleanupRoutine() {
	if mockCache.StartCleanupRoutineFunc != nil {
		mockCache.StartCleanupRoutineFunc()
	}
}

func (mockCache *MockCache) StopCleanupRoutine() {
	if mockCache.StopCleanupRoutineFunc != nil {
		mockCache.StopCleanupRoutineFunc()
	}
}
