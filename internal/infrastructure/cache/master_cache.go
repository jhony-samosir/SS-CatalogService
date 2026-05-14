package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"ss-catalog-service/internal/domain"
	"time"

	"github.com/allegro/bigcache/v3"
)

type masterDataCacheRepository struct {
	cache *bigcache.BigCache
}

func NewMasterDataCacheRepository(eviction time.Duration) (domain.MasterDataCacheRepository, error) {
	config := bigcache.DefaultConfig(eviction)
	config.HardMaxCacheSize = 256 // 256MB for master data is plenty

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize master bigcache: %w", err)
	}

	return &masterDataCacheRepository{cache: cache}, nil
}

func (r *masterDataCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.cache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *masterDataCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// Note: BigCache doesn't support per-key TTL easily, it uses the global LifeWindow.
	// But we can wrap it if needed. For now, global TTL is fine as per config.
	return r.cache.Set(key, data)
}

func (r *masterDataCacheRepository) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(key)
}

func (r *masterDataCacheRepository) InvalidateAll(ctx context.Context) error {
	return r.cache.Reset()
}
