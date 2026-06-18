package cache

import (
	"sync"

	"github.com/LucasM4r/repomind/internal/domain"
)

type SessionCache struct {
	data sync.Map
}

func (c *SessionCache) Get(sessionID string) ([]domain.ChatMessage, bool) {
	val, ok := c.data.Load(sessionID)
	if !ok {
		return nil, false
	}
	return val.([]domain.ChatMessage), true
}

func (c *SessionCache) Set(sessionID string, history []domain.ChatMessage) {
	c.data.Store(sessionID, history)
}
