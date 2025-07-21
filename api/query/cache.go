package query

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

type CacheInterface interface {
	Add(sessionId, userInput, aiOutput string)
	Get(sessionId string) []Context
	StartCleanupRoutine()
	StopCleanupRoutine()
}

type Context struct {
	Role    string
	Content string
}

type ContextCache struct {
	mu              sync.Mutex
	Cache           map[string][]Context
	MaxSize         int
	LastAccessed    map[string]time.Time
	cleanupInterval time.Duration
	sessionTimeout  time.Duration
	stopChan        chan struct{}
}

// maxSize - number of most recent user queries to cache;
// cleanupInterval - how often the session cleanup function runs;
// sessionTimeout - duration after which a session is considered stale;
func NewContextCache(maxSize int, cleanupInterval time.Duration, sessionTimeout time.Duration) (*ContextCache, error) {
	if maxSize < 0 {
		return nil, errors.New("cache size can not be negative")
	}
	if cleanupInterval <= 0 {
		return nil, errors.New("cleanup interval has to be positive")
	}
	if sessionTimeout <= 0 {
		return nil, errors.New("session timeout has to be positive")
	}

	return &ContextCache{
		Cache:           make(map[string][]Context),
		MaxSize:         maxSize,
		LastAccessed:    make(map[string]time.Time),
		cleanupInterval: cleanupInterval,
		sessionTimeout:  sessionTimeout,
		stopChan:        make(chan struct{}),
	}, nil
}

func (cc *ContextCache) Add(sessionId, userInput, aiOutput string) {
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
	cc.LastAccessed[sessionId] = time.Now()
	fmt.Printf("Session %s: Added contexts. Current count: %d\n", sessionId, len(cc.Cache[sessionId]))
}

func (cc *ContextCache) Get(sessionId string) []Context {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.LastAccessed[sessionId] = time.Now()
	fmt.Printf("Session %s: Retrieved contexts.\n", sessionId)
	return cc.Cache[sessionId]
}

func (cc *ContextCache) StartCleanupRoutine() {
	ticker := time.NewTicker(cc.cleanupInterval)
	go func() {
		defer ticker.Stop()
		fmt.Printf("Cleanup routine started, running every %s, session timeout is %s.\n", cc.cleanupInterval, cc.sessionTimeout)
		for {
			select {
			case <-ticker.C:
				cc.cleanupStaleSessions()
			case <-cc.stopChan:
				fmt.Println("Cleanup routine stopped.")
				return
			}
		}
	}()
}

func (cc *ContextCache) StopCleanupRoutine() {
	close(cc.stopChan)
}

func (cc *ContextCache) cleanupStaleSessions() {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	now := time.Now()
	sessionsToDelete := []string{}

	for sessionId, lastAccessed := range cc.LastAccessed {
		if now.Sub(lastAccessed) >= cc.sessionTimeout {
			sessionsToDelete = append(sessionsToDelete, sessionId)
		}
	}

	if len(sessionsToDelete) > 0 {
		fmt.Printf("Cleanup: Found %d expired sessions.\n", len(sessionsToDelete))
	}

	for _, sessionId := range sessionsToDelete {
		delete(cc.Cache, sessionId)
		delete(cc.LastAccessed, sessionId)
		fmt.Printf("Cleanup: Deleted session %s (not accessed for %s).\n", sessionId, now.Sub(cc.LastAccessed[sessionId]))
	}
}

/*
TODO: add tests
*/
