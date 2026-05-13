package product

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
)

type productQueryUsecase struct {
	repo        domain.ProductRepository
	searchRepo  domain.SearchRepository
	cacheRepo   domain.ProductCacheRepository
	sfGroup     *singleflight.Group
	defaultLang string
}

// NewProductQueryUsecase creates a new instance of product query business logic.
func NewProductQueryUsecase(
	repo domain.ProductRepository,
	searchRepo domain.SearchRepository,
	cacheRepo domain.ProductCacheRepository,
	defaultLang string,
) domain.ProductQueryUsecase {
	return &productQueryUsecase{
		repo:        repo,
		searchRepo:  searchRepo,
		cacheRepo:   cacheRepo,
		sfGroup:     &singleflight.Group{},
		defaultLang: defaultLang,
	}
}

func (u *productQueryUsecase) GetAllProducts(ctx context.Context, p domain.Pagination) ([]domain.Product, error) {
	products, err := u.repo.FindAll(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetAllProducts: %w", err)
	}
	return products, nil
}

func (u *productQueryUsecase) GetProductByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	product, err := u.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetProductByPublicID: %w", err)
	}
	if product == nil {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (u *productQueryUsecase) GetProductDetails(ctx context.Context, query domain.GetProductDetailsQuery) (*domain.ProductDetailsResponse, error) {
	langCode := query.LangCode
	if langCode == "" {
		langCode = u.defaultLang
	}

	// 1. Try Cache
	if u.cacheRepo != nil {
		cached, err := u.cacheRepo.GetProductDetails(ctx, query.PublicID, langCode)
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	// 2. Cache Miss - Use Singleflight to prevent stampede
	key := fmt.Sprintf("get_product_details:%s:%s", query.PublicID.String(), langCode)
	val, err, _ := u.sfGroup.Do(key, func() (interface{}, error) {
		product, err := u.repo.GetProductDetails(ctx, query.PublicID, langCode)
		if err != nil {
			return nil, fmt.Errorf("productQueryUsecase.GetProductDetails: %w", err)
		}
		if product == nil {
			return nil, domain.ErrProductNotFound
		}

		// Mapping to DTO
		resp := &domain.ProductDetailsResponse{
			PublicID:  product.PublicID,
			BasePrice: 0,
			Status:    string(product.Status),
		}

		if product.Translation != nil {
			resp.Name = product.Translation.Name
			resp.Description = product.Translation.Description
			resp.ShortDesc = product.Translation.ShortDesc
		} else {
			resp.Name = product.Name
			resp.Description = product.Description
			resp.ShortDesc = product.ShortDesc
		}

		if product.SEO != nil {
			resp.MetaTitle = product.SEO.MetaTitle
			resp.MetaDescription = product.SEO.MetaDescription
		}

		if len(product.Categories) > 0 {
			resp.Categories = make([]string, len(product.Categories))
			for i, c := range product.Categories {
				resp.Categories[i] = c.Name
			}
		}

		if len(product.Tags) > 0 {
			resp.Tags = make([]string, len(product.Tags))
			for i, t := range product.Tags {
				resp.Tags[i] = t.Name
			}
		}

		// 3. Populate Cache
		if u.cacheRepo != nil {
			_ = u.cacheRepo.SetProductDetails(ctx, query.PublicID, langCode, resp)
		}

		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	resp, ok := val.(*domain.ProductDetailsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight for key: %s", key)
	}

	return resp, nil
}

func (u *productQueryUsecase) SearchProducts(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
	// --- Normalization & Validation ---

	// Trim keyword whitespace
	if q.Keyword != nil {
		trimmed := strings.TrimSpace(*q.Keyword)
		q.Keyword = &trimmed
	}

	// Limit normalization (1-100, default 20)
	if q.Limit <= 0 {
		q.Limit = 20
	} else if q.Limit > 100 {
		q.Limit = 100
	}

	// Price range cross-validation
	if q.MinPrice != nil && q.MaxPrice != nil && *q.MinPrice > *q.MaxPrice {
		return nil, fmt.Errorf("%w: min_price cannot be greater than max_price", domain.ErrInvalidInput)
	}

	result, err := u.repo.Search(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.SearchProducts: %w", err)
	}
	return result, nil
}

func (u *productQueryUsecase) FacetedSearch(ctx context.Context, q domain.GetProductSearchQuery) (*domain.FacetedSearchResult, error) {
	if u.searchRepo != nil {
		return u.searchRepo.Search(ctx, q)
	}

	searchResult, err := u.SearchProducts(ctx, q)
	if err != nil {
		return nil, err
	}

	// Stub facets for demonstration
	// In a real implementation, these would come from database aggregations (e.g. Meilisearch facets or SQL GROUP BY)
	facets := []domain.SearchFacet{
		{
			Name: "Category",
			Values: []struct {
				Value string `json:"value"`
				Count int    `json:"count"`
			}{
				{Value: "Snacks", Count: 42},
				{Value: "Drinks", Count: 12},
			},
		},
		{
			Name: "Brand",
			Values: []struct {
				Value string `json:"value"`
				Count int    `json:"count"`
			}{
				{Value: "Sam's", Count: 25},
				{Value: "IndoFood", Count: 18},
			},
		},
	}

	return &domain.FacetedSearchResult{
		Items:      searchResult.Items,
		Facets:     facets,
		TotalHint:  searchResult.TotalHint,
		NextCursor: searchResult.NextCursor,
	}, nil
}
