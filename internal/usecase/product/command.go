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
	outboxRepo domain.OutboxRepository
	txManager  domain.TransactionManager
}

// NewProductCommandUsecase creates a new instance of product command business logic.
func NewProductCommandUsecase(
	repo domain.ProductRepository,
	outboxRepo domain.OutboxRepository,
	txManager domain.TransactionManager,
) domain.ProductCommandUsecase {
	return &productCommandUsecase{
		repo:       repo,
		outboxRepo: outboxRepo,
		txManager:  txManager,
	}
}

func (u *productCommandUsecase) CreateProduct(ctx context.Context, payload domain.CreateProductPayload) (*domain.Product, error) {
	var product *domain.Product

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		product = &domain.Product{
			BaseEntity: domain.BaseEntity{
				PublicID: uuid.New(),
			},
			BrandID: payload.BrandID,
			Name:    payload.Name,
			Slug:    generateSlug(payload.Name),
			Status:  domain.ProductStatusDraft,
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
