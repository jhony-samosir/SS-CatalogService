package product

import (
	"context"
	"errors"
	"testing"

	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)


func TestGetProductDetails(t *testing.T) {
	// ... (existing TestGetProductDetails)
	publicID := uuid.New()
	langCode := "en-US"

	t.Run("Success", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			getProductDetailsFn: func(ctx context.Context, pid uuid.UUID, lang string) (*domain.Product, error) {
				if pid == publicID && lang == langCode {
					return &domain.Product{
						BaseEntity: domain.BaseEntity{PublicID: publicID, ID: 1},
						Name:       "Test Product",
						Translation: &domain.ProductTranslation{
							ProductID: 1,
							Name:      "Test Product EN",
						},
					}, nil
				}
				return nil, nil
			},
		}

		usecase := NewProductQueryUsecase(mockRepo, nil, "id-ID")
		resp, err := usecase.GetProductDetails(context.Background(), domain.GetProductDetailsQuery{
			PublicID: publicID,
			LangCode: langCode,
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if resp == nil {
			t.Fatal("expected response, got nil")
		}

		if resp.Name != "Test Product EN" {
			t.Errorf("expected name 'Test Product EN', got '%s'", resp.Name)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			getProductDetailsFn: func(ctx context.Context, pid uuid.UUID, lang string) (*domain.Product, error) {
				return nil, nil
			},
		}

		usecase := NewProductQueryUsecase(mockRepo, nil, "id-ID")
		_, err := usecase.GetProductDetails(context.Background(), domain.GetProductDetailsQuery{
			PublicID: publicID,
			LangCode: langCode,
		})

		if err != domain.ErrProductNotFound {
			t.Errorf("expected ErrProductNotFound, got %v", err)
		}
	})
}

func TestSearchProducts(t *testing.T) {
	t.Run("Success_NoFilters", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			searchFn: func(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
				return &domain.ProductSearchResult{
					Items: []domain.Product{{Name: "Product 1"}},
				}, nil
			},
		}

		usecase := NewProductQueryUsecase(mockRepo, nil, "en")
		res, err := usecase.SearchProducts(context.Background(), domain.GetProductSearchQuery{})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(res.Items) != 1 {
			t.Errorf("expected 1 item, got %d", len(res.Items))
		}
	})

	t.Run("Keyword_Trimming", func(t *testing.T) {
		keyword := "  shoes  "
		mockRepo := &mockProductRepository{
			searchFn: func(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
				if *q.Keyword != "shoes" {
					t.Errorf("expected keyword 'shoes', got '%s'", *q.Keyword)
				}
				return &domain.ProductSearchResult{}, nil
			},
		}

		usecase := NewProductQueryUsecase(mockRepo, nil, "en")
		_, _ = usecase.SearchProducts(context.Background(), domain.GetProductSearchQuery{Keyword: &keyword})
	})

	t.Run("PriceRange_Validation_Error", func(t *testing.T) {
		min := 100.0
		max := 50.0
		mockRepo := &mockProductRepository{}

		usecase := NewProductQueryUsecase(mockRepo, nil, "en")
		_, err := usecase.SearchProducts(context.Background(), domain.GetProductSearchQuery{
			MinPrice: &min,
			MaxPrice: &max,
		})

		if err == nil {
			t.Error("expected validation error for min > max price, got nil")
		} else if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("Limit_Normalization", func(t *testing.T) {
		mockRepo := &mockProductRepository{
			searchFn: func(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
				if q.Limit != 100 {
					t.Errorf("expected limit 100, got %d", q.Limit)
				}
				return &domain.ProductSearchResult{}, nil
			},
		}

		usecase := NewProductQueryUsecase(mockRepo, nil, "en")
		_, _ = usecase.SearchProducts(context.Background(), domain.GetProductSearchQuery{Limit: 999})
	})
}
