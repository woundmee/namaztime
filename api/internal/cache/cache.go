package cache

import (
	"log/slog"
	"namaztimeApi/models"
	"sync"
)

type Cache struct {
	logger *slog.Logger
	mu     sync.RWMutex
	data   []models.NamazTime
}

func New(logger *slog.Logger) *Cache {
	return &Cache{
		logger: logger,
		data:   []models.NamazTime{},
		mu:     sync.RWMutex{},
	}
}

func (c *Cache) Set(data []models.NamazTime) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// c.data = data
	c.data = append([]models.NamazTime(nil), data...)
}

func (c *Cache) Get() ([]models.NamazTime, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.data) == 0 {
		c.logger.Warn("кэш пустой")
		return nil, false
	}

	return c.data, true
}
