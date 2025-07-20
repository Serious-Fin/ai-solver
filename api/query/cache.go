package query

import (
	"errors"
	"sync"
	"time"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type Context struct {
	Role    string
	Content string
}

type ContextCache struct {
	mu      sync.Mutex
	Cache   map[string][]Context
	MaxSize int
}

type TimeRecord struct {
	LastModified time.Time
	SessionId    string
}

// maxSize - number of most recent user queries to cache
func NewContextCache(maxSize int) (*ContextCache, error) {
	if maxSize < 0 {
		return nil, errors.New("cache size can not be negative")
	}

	return &ContextCache{
		Cache:   make(map[string][]Context),
		MaxSize: maxSize,
	}, nil
}

func (cc *ContextCache) Add(sessionId string, userInput string, aiOutput string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cache := cc.Cache[sessionId]
	cache = append(cache, Context{
		Role:    RoleUser,
		Content: userInput,
	})
	cache = append(cache, Context{
		Role:    RoleAssistant,
		Content: aiOutput,
	})

	if len(cache) > cc.MaxSize*2 {
		valuesToRemove := (len(cache) - (cc.MaxSize * 2))
		cache = cache[valuesToRemove:]
	}
	cc.Cache[sessionId] = cache
}

func (cc *ContextCache) Get(sessionId string) []Context {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.Cache[sessionId]
}

/*
TODO: Implement automatic cache expiration mechanism
TODO: add tests
*/
