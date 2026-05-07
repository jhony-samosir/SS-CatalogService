package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ss-catalog-service/internal/domain"
	"ss-catalog-service/internal/mocks"
	"ss-catalog-service/internal/usecase/product"
)

func TestGetProductDetails(t *testing.T) {
	publicID := uuid.New()
	langCode := "en-US"

	testCases := []struct {
		name          string
		query         domain.GetProductDetailsQuery
		setupMocks    func(repo *mocks.MockProductRepository)
		expectedError error
		expectedName  string
	}{
		{
			name: "Success",
			query: domain.GetProductDetailsQuery{
				PublicID: publicID,
				LangCode: langCode,
			},
			setupMocks: func(repo *mocks.MockProductRepository) {
				repo.On("GetProductDetails", mock.Anything, publicID, langCode).Return(&domain.Product{
					BaseEntity: domain.BaseEntity{PublicID: publicID, ID: 1},
					Name:       "Test Product",
					Translation: &domain.ProductTranslation{
						ProductID: 1,
						Name:      "Test Product EN",
					},
				}, nil).Once()
			},
			expectedError: nil,
			expectedName:  "Test Product EN",
		},
		{
			name: "NotFound",
			query: domain.GetProductDetailsQuery{
				PublicID: publicID,
				LangCode: langCode,
			},
			setupMocks: func(repo *mocks.MockProductRepository) {
				repo.On("GetProductDetails", mock.Anything, publicID, langCode).Return(nil, nil).Once()
			},
			expectedError: domain.ErrProductNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProductRepository)
			tc.setupMocks(mockRepo)

			usecase := product.NewProductQueryUsecase(mockRepo, nil, "id-ID")
			resp, err := usecase.GetProductDetails(context.Background(), tc.query)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if resp != nil {
					assert.Equal(t, tc.expectedName, resp.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSearchProducts(t *testing.T) {
	testCases := []struct {
		name          string
		query         domain.GetProductSearchQuery
		setupMocks    func(repo *mocks.MockProductRepository)
		expectedError error
		assertResult  func(t *testing.T, res *domain.ProductSearchResult)
	}{
		{
			name:  "Success_NoFilters",
			query: domain.GetProductSearchQuery{},
			setupMocks: func(repo *mocks.MockProductRepository) {
				repo.On("Search", mock.Anything, mock.MatchedBy(func(q domain.GetProductSearchQuery) bool {
					return true
				})).Return(&domain.ProductSearchResult{
					Items: []domain.Product{{Name: "Product 1"}},
				}, nil).Once()
			},
			expectedError: nil,
			assertResult: func(t *testing.T, res *domain.ProductSearchResult) {
				assert.Len(t, res.Items, 1)
			},
		},
		{
			name: "Keyword_Trimming",
			query: func() domain.GetProductSearchQuery {
				kw := "  shoes  "
				return domain.GetProductSearchQuery{Keyword: &kw}
			}(),
			setupMocks: func(repo *mocks.MockProductRepository) {
				repo.On("Search", mock.Anything, mock.MatchedBy(func(q domain.GetProductSearchQuery) bool {
					return *q.Keyword == "shoes"
				})).Return(&domain.ProductSearchResult{}, nil).Once()
			},
			expectedError: nil,
			assertResult: func(t *testing.T, res *domain.ProductSearchResult) {
				assert.NotNil(t, res)
			},
		},
		{
			name: "PriceRange_Validation_Error",
			query: func() domain.GetProductSearchQuery {
				min := 100.0
				max := 50.0
				return domain.GetProductSearchQuery{MinPrice: &min, MaxPrice: &max}
			}(),
			setupMocks: func(repo *mocks.MockProductRepository) {
				// Should fail before reaching DB
			},
			expectedError: domain.ErrInvalidInput,
		},
		{
			name: "Limit_Normalization",
			query: domain.GetProductSearchQuery{Limit: 999},
			setupMocks: func(repo *mocks.MockProductRepository) {
				repo.On("Search", mock.Anything, mock.MatchedBy(func(q domain.GetProductSearchQuery) bool {
					return q.Limit == 100
				})).Return(&domain.ProductSearchResult{}, nil).Once()
			},
			expectedError: nil,
			assertResult: func(t *testing.T, res *domain.ProductSearchResult) {
				assert.NotNil(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProductRepository)
			tc.setupMocks(mockRepo)

			usecase := product.NewProductQueryUsecase(mockRepo, nil, "en")
			res, err := usecase.SearchProducts(context.Background(), tc.query)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))
			} else {
				assert.NoError(t, err)
				if tc.assertResult != nil {
					tc.assertResult(t, res)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
