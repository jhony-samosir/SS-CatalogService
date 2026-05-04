package product

import (
	"context"
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
	repo domain.ProductRepository
}

// NewProductCommandUsecase creates a new instance of product command business logic.
func NewProductCommandUsecase(repo domain.ProductRepository) domain.ProductCommandUsecase {
	return &productCommandUsecase{repo: repo}
}

func (u *productCommandUsecase) CreateProduct(ctx context.Context, payload domain.CreateProductPayload) (*domain.Product, error) {
	product := &domain.Product{
		BaseEntity: domain.BaseEntity{
			PublicID: uuid.New(),
		},
		BrandID: payload.BrandID,
		Name:    payload.Name,
		Slug:    generateSlug(payload.Name),
		Status:  domain.ProductStatusDraft,
	}

	if err := u.repo.Create(ctx, product); err != nil {
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
