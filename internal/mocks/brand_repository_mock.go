package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"ss-catalog-service/internal/domain"
)

type MockBrandRepository struct {
	mock.Mock
}

func (m *MockBrandRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Brand, error) {
	args := m.Called(ctx, p)
	return args.Get(0).([]domain.Brand), args.Error(1)
}

func (m *MockBrandRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Brand, error) {
	args := m.Called(ctx, publicID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Brand), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBrandRepository) Create(ctx context.Context, brand *domain.Brand) error {
	args := m.Called(ctx, brand)
	return args.Error(0)
}

func (m *MockBrandRepository) Update(ctx context.Context, brand *domain.Brand) error {
	args := m.Called(ctx, brand)
	return args.Error(0)
}

func (m *MockBrandRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	args := m.Called(ctx, publicID)
	return args.Error(0)
}

func (m *MockBrandRepository) CountProducts(ctx context.Context, brandID int) (int64, error) {
	args := m.Called(ctx, brandID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockBrandRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}
