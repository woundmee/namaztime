package cache

import (
	"log/slog"
	"sync"
)

type Cache struct {
	logger *slog.Logger
	mu     sync.RWMutex
	data   []byte
}

func New(logger *slog.Logger) *Cache {
	return &Cache{
		logger: logger,
		data:   []byte{},
		mu:     sync.RWMutex{},
	}
}

func (c *Cache) Set(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = data
}

func (c *Cache) Get() ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.data) == 0 {
		c.logger.Warn("кэш пустой")
		return nil, false
	}

	return c.data, true
}
