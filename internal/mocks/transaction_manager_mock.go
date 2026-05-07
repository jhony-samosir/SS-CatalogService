package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockTransactionManager is a mock implementation of domain.TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if f, ok := args.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		return f(ctx, fn)
	}
	return args.Error(0)
}
