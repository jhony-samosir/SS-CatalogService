package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"ss-catalog-service/internal/domain"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Category, error) {
	args := m.Called(ctx, p)
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByID(ctx context.Context, id int) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Category, error) {
	args := m.Called(ctx, publicID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Category), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	args := m.Called(ctx, publicID)
	return args.Error(0)
}

func (m *MockCategoryRepository) CountChildren(ctx context.Context, parentID int) (int64, error) {
	args := m.Called(ctx, parentID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockCategoryRepository) CountProducts(ctx context.Context, categoryID int) (int64, error) {
	args := m.Called(ctx, categoryID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockCategoryRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}
