package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
)

type MemoryCache struct {
	cache *bigcache.BigCache
}

func NewMemoryCache(duration time.Duration) (*MemoryCache, error) {
	config := bigcache.DefaultConfig(duration)
	config.HardMaxCacheSize = 1024 // 1GB limit for safety

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &MemoryCache{cache: cache}, nil
}

func (m *MemoryCache) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.cache.Set(key, data)
}

func (m *MemoryCache) Get(key string, dest interface{}) error {
	data, err := m.cache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (m *MemoryCache) Delete(key string) error {
	return m.cache.Delete(key)
}

func (m *MemoryCache) Reset() error {
	return m.cache.Reset()
}
