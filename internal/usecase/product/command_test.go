package product

import (
	"context"
	"testing"

	"ss-catalog-service/internal/domain"
)

type mockProductCmdRepository struct {
	domain.ProductRepository
	createFn func(ctx context.Context, p *domain.Product) error
}

func (m *mockProductCmdRepository) Create(ctx context.Context, p *domain.Product) error {
	return m.createFn(ctx, p)
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

	usecase := NewProductCommandUsecase(mockProductRepo, mockOutboxRepo, mockTransactionManager)

	brandID := 1
	payload := domain.CreateProductPayload{
		Name:    "Test Product",
		BrandID: &brandID,
	}

	product, err := usecase.CreateProduct(context.Background(), payload)

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
