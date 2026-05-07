package product

import (
	"context"
	"testing"

	"ss-catalog-service/internal/domain"
	"errors"
	"github.com/google/uuid"
)

type mockProductCmdRepository struct {
	domain.ProductRepository
	createFn         func(ctx context.Context, p *domain.Product) error
	updateFn         func(ctx context.Context, p *domain.Product) error
	findByPublicIDFn func(ctx context.Context, pid uuid.UUID) (*domain.Product, error)
}

func (m *mockProductCmdRepository) Create(ctx context.Context, p *domain.Product) error {
	return m.createFn(ctx, p)
}

func (m *mockProductCmdRepository) Update(ctx context.Context, p *domain.Product) error {
	return m.updateFn(ctx, p)
}

type mockProductRepository struct {
	domain.ProductRepository
	findByPublicIDFn    func(ctx context.Context, pid uuid.UUID) (*domain.Product, error)
	getProductDetailsFn func(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.Product, error)
	searchFn            func(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error)
}

func (m *mockProductRepository) FindByPublicID(ctx context.Context, pid uuid.UUID) (*domain.Product, error) {
	return m.findByPublicIDFn(ctx, pid)
}

func (m *mockProductRepository) GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.Product, error) {
	return m.getProductDetailsFn(ctx, publicID, langCode)
}

func (m *mockProductRepository) Search(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
	return m.searchFn(ctx, q)
}

type mockOutboxCmdRepository struct {
	domain.OutboxRepository
	saveFn func(ctx context.Context, e *domain.OutboxEvent) error
}

func (m *mockOutboxCmdRepository) Save(ctx context.Context, e *domain.OutboxEvent) error {
	return m.saveFn(ctx, e)
}

type mockTxManager struct {
	domain.TransactionManager
}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCreateProductWithOutbox(t *testing.T) {
	mockProductRepo := &mockProductCmdRepository{
		createFn: func(ctx context.Context, p *domain.Product) error {
			p.ID = 1
			return nil
		},
	}

	outboxSaved := false
	mockOutboxRepo := &mockOutboxCmdRepository{
		saveFn: func(ctx context.Context, e *domain.OutboxEvent) error {
			if e.EventType == "ProductCreated" && e.AggregateID == 1 {
				outboxSaved = true
			}
			return nil
		},
	}

	mockTransactionManager := &mockTxManager{}

	usecase := NewProductCommandUsecase(mockProductRepo, nil, mockOutboxRepo, mockTransactionManager)

	brandID := 1
	payload := domain.CreateProductPayload{
		Name:    "Test Product",
		BrandID: &brandID,
	}

	sellerID := 123
	ctx := domain.ContextWithUser(context.Background(), domain.UserContext{
		SellerID: &sellerID,
	})

	product, err := usecase.CreateProduct(ctx, payload)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if !outboxSaved {
		t.Error("expected outbox event to be saved")
	}
}

func TestUpdateProductAuthorization(t *testing.T) {
	publicID := uuid.New()
	sellerID := 123
	otherSellerID := 456

	mockRepo := &mockProductCmdRepository{
		ProductRepository: &mockProductRepository{
			findByPublicIDFn: func(ctx context.Context, pid uuid.UUID) (*domain.Product, error) {
				if pid == publicID {
					return &domain.Product{
						BaseEntity: domain.BaseEntity{PublicID: publicID},
						SellerID:   &sellerID,
					}, nil
				}
				return nil, nil
			},
		},
		updateFn: func(ctx context.Context, p *domain.Product) error {
			return nil
		},
	}

	mockTransactionManager := &mockTxManager{}
	mockOutboxRepo := &mockOutboxCmdRepository{
		saveFn: func(ctx context.Context, e *domain.OutboxEvent) error {
			return nil
		},
	}
	usecase := NewProductCommandUsecase(mockRepo, nil, mockOutboxRepo, mockTransactionManager)

	t.Run("Success_OwnerMatches", func(t *testing.T) {
		ctx := domain.ContextWithUser(context.Background(), domain.UserContext{
			SellerID: &sellerID,
		})

		err := usecase.UpdateProduct(ctx, domain.UpdateProductPayload{
			PublicID: publicID,
			Name:     "New Name",
		})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Error_OwnerMismatches", func(t *testing.T) {
		ctx := domain.ContextWithUser(context.Background(), domain.UserContext{
			SellerID: &otherSellerID,
		})

		err := usecase.UpdateProduct(ctx, domain.UpdateProductPayload{
			PublicID: publicID,
			Name:     "New Name",
		})

		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("Error_NoUserInContext", func(t *testing.T) {
		err := usecase.UpdateProduct(context.Background(), domain.UpdateProductPayload{
			PublicID: publicID,
			Name:     "New Name",
		})

		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("Success_AdminBypass", func(t *testing.T) {
		ctx := domain.ContextWithUser(context.Background(), domain.UserContext{
			Roles: []string{"admin"},
		})

		err := usecase.UpdateProduct(ctx, domain.UpdateProductPayload{
			PublicID: publicID,
			Name:     "Admin Update",
		})

		if err != nil {
			t.Errorf("expected no error for admin, got %v", err)
		}
	})
}
