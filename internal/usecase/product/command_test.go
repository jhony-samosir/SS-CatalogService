package product_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/internal/mocks"
	"ss-catalog-service/internal/usecase/product"
)

func TestCreateProductCommand(t *testing.T) {
	sellerID := 123
	ctx := domain.ContextWithUser(context.Background(), domain.UserContext{
		SellerID: &sellerID,
	})

	testCases := []struct {
		name          string
		payload       domain.CreateProductPayload
		setupMocks    func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager)
		expectedError error
	}{
		{
			name: "Success: Product created",
			payload: domain.CreateProductPayload{
				Name: "New Product",
			},
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {
				tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Once()

				// State Verification: check if name and seller_id match
				repo.On("Create", mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
					return p.Name == "New Product" && *p.SellerID == sellerID
				})).Return(nil).Once()

				outbox.On("Save", mock.Anything, mock.AnythingOfType("*domain.OutboxEvent")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Failure: Empty product name validation",
			payload: domain.CreateProductPayload{
				Name: "",
			},
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {},
			expectedError: domain.ErrInvalidProductName,
		},
		{
			name: "Failure: Repository database error",
			payload: domain.CreateProductPayload{
				Name: "New Product",
			},
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {
				tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Once()

				repo.On("Create", mock.Anything, mock.Anything).Return(domain.ErrInternalDatabase).Once()
			},
			expectedError: domain.ErrInternalDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProductRepository)
			mockOutbox := new(mocks.MockOutboxRepository)
			mockTx := new(mocks.MockTransactionManager)

			tc.setupMocks(mockRepo, mockOutbox, mockTx)
			usecase := product.NewProductCommandUsecase(mockRepo, nil, mockOutbox, mockTx)

			_, err := usecase.CreateProduct(ctx, tc.payload)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockOutbox.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}

func TestUpdateProductAuthorization(t *testing.T) {
	publicID := uuid.New()
	sellerID := 123
	otherSellerID := 456

	basePayload := domain.UpdateProductPayload{
		PublicID: publicID,
		Name:     "Updated Name",
	}

	testCases := []struct {
		name          string
		userCtx       domain.UserContext
		setupMocks    func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager)
		expectedError error
	}{
		{
			name: "Success: Owner Matches",
			userCtx: domain.UserContext{
				SellerID: &sellerID,
			},
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {
				tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Once()

				repo.On("FindByPublicID", mock.Anything, publicID).Return(&domain.Product{
					BaseEntity: domain.BaseEntity{PublicID: publicID, ID: 1},
					SellerID:   &sellerID,
				}, nil).Once()

				// State Verification: ensure only updated fields are changed
				repo.On("Update", mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
					return p.Name == "Updated Name"
				})).Return(nil).Once()
				
				outbox.On("Save", mock.Anything, mock.AnythingOfType("*domain.OutboxEvent")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Success: Admin Bypass",
			userCtx: domain.UserContext{
				Roles: []string{"admin"},
			},
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {
				tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Once()

				repo.On("FindByPublicID", mock.Anything, publicID).Return(&domain.Product{
					BaseEntity: domain.BaseEntity{PublicID: publicID, ID: 1},
					SellerID:   &sellerID, 
				}, nil).Once()

				repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				outbox.On("Save", mock.Anything, mock.AnythingOfType("*domain.OutboxEvent")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Error: Owner Mismatches",
			userCtx: domain.UserContext{
				SellerID: &otherSellerID,
			},
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {
				tx.On("WithTransaction", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Once()

				repo.On("FindByPublicID", mock.Anything, publicID).Return(&domain.Product{
					BaseEntity: domain.BaseEntity{PublicID: publicID, ID: 1},
					SellerID:   &sellerID,
				}, nil).Once()
			},
			expectedError: domain.ErrUnauthorized,
		},
		{
			name: "Error: No User in Context",
			userCtx: domain.UserContext{}, 
			setupMocks: func(repo *mocks.MockProductRepository, outbox *mocks.MockOutboxRepository, tx *mocks.MockTransactionManager) {},
			expectedError: domain.ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProductRepository)
			mockOutbox := new(mocks.MockOutboxRepository)
			mockTx := new(mocks.MockTransactionManager)

			tc.setupMocks(mockRepo, mockOutbox, mockTx)
			usecase := product.NewProductCommandUsecase(mockRepo, nil, mockOutbox, mockTx)

			var ctx context.Context
			if tc.name == "Error: No User in Context" {
				ctx = context.Background()
			} else {
				ctx = domain.ContextWithUser(context.Background(), tc.userCtx)
			}

			err := usecase.UpdateProduct(ctx, basePayload)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockOutbox.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}
