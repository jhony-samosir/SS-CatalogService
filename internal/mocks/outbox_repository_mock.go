package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"ss-catalog-service/internal/domain"
)

// MockOutboxRepository is a mock implementation of domain.OutboxRepository
type MockOutboxRepository struct {
	mock.Mock
}

func (m *MockOutboxRepository) Save(ctx context.Context, event *domain.OutboxEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockOutboxRepository) FetchPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.OutboxEvent), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOutboxRepository) MarkAsPublished(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOutboxRepository) MarkAsFailed(ctx context.Context, id int, errorMessage string) error {
	args := m.Called(ctx, id, errorMessage)
	return args.Error(0)
}
