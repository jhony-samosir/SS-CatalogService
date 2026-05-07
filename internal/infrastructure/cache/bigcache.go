package cache

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

type productCacheRepository struct {
	cache       *bigcache.BigCache
	activeLangs []string // Supported languages for invalidation
}

// NewProductCacheRepository creates a new instance of BigCache-based product cache.
func NewProductCacheRepository(eviction time.Duration, activeLangs []string) (domain.ProductCacheRepository, error) {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         eviction,
		CleanWindow:        1 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            false,
		HardMaxCacheSize:   512, // MB
	}

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bigcache: %w", err)
	}

	return &productCacheRepository{
		cache:       cache,
		activeLangs: activeLangs,
	}, nil
}

func (r *productCacheRepository) GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.ProductDetailsResponse, error) {
	key := r.generateKey(publicID, langCode)
	data, err := r.cache.Get(key)
	if err != nil {
		return nil, err
	}

	var product domain.ProductDetailsResponse
	// Use MsgPack for better performance and zero-GC efficiency
	if err := msgpack.Unmarshal(data, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached product: %w", err)
	}

	return &product, nil
}

func (r *productCacheRepository) SetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string, product *domain.ProductDetailsResponse) error {
	key := r.generateKey(publicID, langCode)
	data, err := msgpack.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product for cache: %w", err)
	}

	return r.cache.Set(key, data)
}

func (r *productCacheRepository) InvalidateProductDetails(ctx context.Context, publicID uuid.UUID) error {
	// Explicit Deletion Strategy: Delete for all known languages
	for _, lang := range r.activeLangs {
		key := r.generateKey(publicID, lang)
		_ = r.cache.Delete(key)
	}
	return nil
}

func (r *productCacheRepository) generateKey(publicID uuid.UUID, langCode string) string {
	return fmt.Sprintf("prod_details:%s:%s", publicID.String(), langCode)
}
