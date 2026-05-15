package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"ss-catalog-service/internal/domain"
)

// MockProductRepository is a mock implementation of domain.ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Product, int64, error) {
	args := m.Called(ctx, p)
	return args.Get(0).([]domain.Product), int64(args.Int(1)), args.Error(2)
}

func (m *MockProductRepository) FindByID(ctx context.Context, id int) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, publicID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) GetProductDetails(ctx context.Context, publicID uuid.UUID, langCode string) (*domain.Product, error) {
	args := m.Called(ctx, publicID, langCode)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Search(ctx context.Context, q domain.GetProductSearchQuery) (*domain.ProductSearchResult, error) {
	args := m.Called(ctx, q)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.ProductSearchResult), args.Error(1)
	}
	return nil, args.Error(1)
}
