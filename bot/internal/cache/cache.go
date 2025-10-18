package cache

import (
	"log/slog"
	"sync"
	"time"
)

type Cache struct {
	logger *slog.Logger
	mu     sync.RWMutex
	data   []byte
	//ttl time.Duration
	// expiresAt time.Time
}

func New(logger *slog.Logger) *Cache {
	return &Cache{
		logger: logger,
		data:   []byte{},
		mu:     sync.RWMutex{},
	}
}

func (c *Cache) Set(data []byte) {
	if c == nil {
		panic("🔴 cache is nil!")
		// return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = data
}

func (c *Cache) Get() ([]byte, bool) {

	if c == nil {
		panic("🔴 cache is nil!")
		// return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.data) == 0 {
		c.logger.Warn("кэш пустой")
		return nil, false
	}

	return c.data, true
}

// from utc+7
func (c *Cache) CalculateMidnightUtc7() time.Time {
	loc := time.FixedZone("UTC+7", 7*60*60)
	now := time.Now().In(loc)

	return time.Date(
		now.Year(), now.Month(), now.Day()+1,
		0, 0, 0, 0, loc,
	)
}
