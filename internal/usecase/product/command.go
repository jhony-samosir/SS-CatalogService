package product

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type productCommandUsecase struct {
	repo       domain.ProductRepository
	cacheRepo  domain.ProductCacheRepository
	outboxRepo domain.OutboxRepository
	txManager  domain.TransactionManager
}

// NewProductCommandUsecase creates a new instance of product command business logic.
func NewProductCommandUsecase(
	repo domain.ProductRepository,
	cacheRepo domain.ProductCacheRepository,
	outboxRepo domain.OutboxRepository,
	txManager domain.TransactionManager,
) domain.ProductCommandUsecase {
	return &productCommandUsecase{
		repo:       repo,
		cacheRepo:  cacheRepo,
		outboxRepo: outboxRepo,
		txManager:  txManager,
	}
}

func (u *productCommandUsecase) CreateProduct(ctx context.Context, payload domain.CreateProductPayload) (*domain.Product, error) {
	// Extract User Context for SellerID binding
	userCtx, ok := domain.UserFromContext(ctx)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	// Validation
	if payload.Name == "" {
		return nil, domain.ErrInvalidProductName
	}

	var product *domain.Product

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		product = &domain.Product{
			BaseEntity: domain.BaseEntity{
				PublicID: uuid.New(),
			},
			BrandID:  payload.BrandID,
			SellerID: userCtx.SellerID, // Bind SellerID from context
			Name:     payload.Name,
			Slug:     generateSlug(payload.Name),
			Status:   domain.ProductStatusDraft,
		}

		// 1. Save Product
		if err := u.repo.Create(txCtx, product); err != nil {
			return fmt.Errorf("failed to save product: %w", err)
		}

		// 2. Prepare Outbox Event
		eventPayload, _ := json.Marshal(map[string]interface{}{
			"id":        product.ID,
			"public_id": product.PublicID,
			"name":      product.Name,
			"slug":      product.Slug,
		})

		outboxEvent := &domain.OutboxEvent{
			EventType:     "ProductCreated",
			AggregateType: "Product",
			AggregateID:   product.ID,
			Payload:       eventPayload,
			Status:        domain.OutboxStatusPending,
		}

		// 3. Save Outbox Event (in the same transaction)
		if err := u.outboxRepo.Save(txCtx, outboxEvent); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("productCommandUsecase.CreateProduct: %w", err)
	}

	return product, nil
}

func (u *productCommandUsecase) UpdateProduct(ctx context.Context, payload domain.UpdateProductPayload) error {
	// 1. Extract User Context
	userCtx, ok := domain.UserFromContext(ctx)
	if !ok {
		return domain.ErrUnauthorized
	}

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 2. Fetch the existing product to verify ownership
		product, err := u.repo.FindByPublicID(txCtx, payload.PublicID)
		if err != nil {
			return err
		}
		if product == nil {
			return domain.ErrProductNotFound
		}

		// 3. Strict Authorization Check (IDOR Prevention + RBAC Bypass)
		isAdmin := false
		for _, role := range userCtx.Roles {
			if role == "admin" {
				isAdmin = true
				break
			}
		}

		isOwner := product.SellerID != nil && userCtx.SellerID != nil && *product.SellerID == *userCtx.SellerID

		if !isOwner && !isAdmin {
			return domain.ErrUnauthorized
		}

		// 4. Update fields
		product.Name = payload.Name
		product.Description = payload.Description
		product.Status = payload.Status

		// 5. Save changes
		if err := u.repo.Update(txCtx, product); err != nil {
			return fmt.Errorf("failed to update product: %w", err)
		}

		// 6. Prepare Outbox Event (Consistency)
		eventPayload, _ := json.Marshal(map[string]interface{}{
			"id":        product.ID,
			"public_id": product.PublicID,
			"status":    product.Status,
		})

		outboxEvent := &domain.OutboxEvent{
			EventType:     "ProductUpdated",
			AggregateType: "Product",
			AggregateID:   product.ID,
			Payload:       eventPayload,
			Status:        domain.OutboxStatusPending,
		}

		if err := u.outboxRepo.Save(txCtx, outboxEvent); err != nil {
			return fmt.Errorf("failed to save outbox event: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 7. Invalidate Cache
	if u.cacheRepo != nil {
		_ = u.cacheRepo.InvalidateProductDetails(ctx, payload.PublicID)
	}

	return nil
}

// generateSlug converts a product name to a URL-safe slug.
// e.g. "Baju Batik Élégant!" -> "baju-batik-elegant"
func generateSlug(name string) string {
	// Normalize unicode characters (e.g. é -> e)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Remove diacritical marks
	}), norm.NFC)
	result, _, err := transform.String(t, name)
	if err != nil {
		result = name
	}

	// Lowercase
	result = strings.ToLower(result)

	// Replace non-alphanumeric characters with hyphen
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	result = reg.ReplaceAllString(result, "-")

	// Trim leading/trailing hyphens
	result = strings.Trim(result, "-")

	return result
}
