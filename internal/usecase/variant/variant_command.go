package variant

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type variantCommandUsecase struct {
	variantRepo domain.VariantRepository
	productRepo domain.ProductRepository
	txManager   domain.TransactionManager
}

// NewVariantCommandUsecase creates a new instance of variant command business logic.
func NewVariantCommandUsecase(
	repo domain.VariantRepository,
	productRepo domain.ProductRepository,
	txManager domain.TransactionManager,
) domain.VariantCommandUsecase {
	return &variantCommandUsecase{
		variantRepo: repo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}

func (u *variantCommandUsecase) CreateProductVariant(ctx context.Context, payload domain.CreateVariantPayload) (*domain.ProductVariant, error) {
	// 1. Business Validation: Ensure Product (SPU) exists
	spu, err := u.productRepo.FindByID(ctx, payload.ProductID)
	if err != nil {
		return nil, fmt.Errorf("variantCommandUsecase.CreateProductVariant: %w", err)
	}
	if spu == nil {
		return nil, domain.ErrProductNotFound
	}

	// 2. Prepare Domain Entity
	variant := &domain.ProductVariant{
		BaseEntity: domain.BaseEntity{
			PublicID: uuid.New(),
		},
		ProductID:  payload.ProductID,
		SKU:        payload.SKU,
		Barcode:    payload.Barcode,
		Name:       payload.Name,
		IsDefault:  payload.IsDefault,
		IsActive:   true,
		WeightGram: payload.WeightGram,
	}

	// 2. Orchestrate Transaction
	err = u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// A. Save Variant
		if err := u.variantRepo.CreateVariant(txCtx, variant); err != nil {
			return fmt.Errorf("failed to create variant: %w", err)
		}

		// B. Link Attributes
		if len(payload.Attributes) > 0 {
			attrs := make([]domain.ProductVariantAttribute, len(payload.Attributes))
			for i, attr := range payload.Attributes {
				attrs[i] = domain.ProductVariantAttribute{
					BaseEntity: domain.BaseEntity{
						PublicID: uuid.New(),
					},
					VariantID:        variant.ID,
					AttributeID:      attr.AttributeID,
					AttributeValueID: attr.AttributeValueID,
				}
			}
			if err := u.variantRepo.CreateVariantAttributes(txCtx, attrs); err != nil {
				return fmt.Errorf("failed to link variant attributes: %w", err)
			}
		}

		// C. Save Images
		if len(payload.Images) > 0 {
			images := make([]domain.ProductImage, len(payload.Images))
			for i, img := range payload.Images {
				images[i] = domain.ProductImage{
					BaseEntity: domain.BaseEntity{
						PublicID: uuid.New(),
					},
					ProductID: &payload.ProductID,
					VariantID: &variant.ID,
					URL:       img.URL,
					AltText:   img.AltText,
					SortOrder: i,
					IsPrimary: img.IsPrimary,
				}
			}
			if err := u.variantRepo.CreateVariantImages(txCtx, images); err != nil {
				return fmt.Errorf("failed to save variant images: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("CreateProductVariant transaction failed: %w", err)
	}

	return variant, nil
}
