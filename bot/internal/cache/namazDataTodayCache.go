package cache

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
