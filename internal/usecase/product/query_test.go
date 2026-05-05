package product

import (
	"context"
	"testing"

	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type mockProductRepository struct {
	domain.ProductRepository
	getProductDetailsFn func(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.Product, error)
}

func (m *mockProductRepository) GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.Product, error) {
	return m.getProductDetailsFn(ctx, publicID, langCode)
}

func TestGetProductDetails(t *testing.T) {
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

		usecase := NewProductQueryUsecase(mockRepo, "id-ID")
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

		usecase := NewProductQueryUsecase(mockRepo, "id-ID")
		_, err := usecase.GetProductDetails(context.Background(), domain.GetProductDetailsQuery{
			PublicID: publicID,
			LangCode: langCode,
		})

		if err != domain.ErrProductNotFound {
			t.Errorf("expected ErrProductNotFound, got %v", err)
		}
	})
}
