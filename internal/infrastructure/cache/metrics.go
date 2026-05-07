package cache

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "catalog_cache_hits_total",
		Help: "Total number of cache hits",
	}, []string{"cache_name"})

	cacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "catalog_cache_misses_total",
		Help: "Total number of cache misses",
	}, []string{"cache_name"})
)

type productCacheMetricsDecorator struct {
	base domain.ProductCacheRepository
}

// NewProductCacheMetricsDecorator wraps a cache repository with Prometheus metrics.
func NewProductCacheMetricsDecorator(base domain.ProductCacheRepository) domain.ProductCacheRepository {
	return &productCacheMetricsDecorator{base: base}
}

func (d *productCacheMetricsDecorator) GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.ProductDetailsResponse, error) {
	resp, err := d.base.GetProductDetails(ctx, publicID, langCode)
	if err != nil {
		cacheMisses.WithLabelValues("product_details").Inc()
		return nil, err
	}
	cacheHits.WithLabelValues("product_details").Inc()
	return resp, nil
}

func (d *productCacheMetricsDecorator) SetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string, product *domain.ProductDetailsResponse) error {
	return d.base.SetProductDetails(ctx, publicID, langCode, product)
}

func (d *productCacheMetricsDecorator) InvalidateProductDetails(ctx context.Context, publicID uuid.UUID) error {
	return d.base.InvalidateProductDetails(ctx, publicID)
}
