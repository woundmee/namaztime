package cache

import (
	"log/slog"
	"sync"
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
	}
}
