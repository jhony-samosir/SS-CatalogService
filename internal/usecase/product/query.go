package product

import (
	"context"
	"fmt"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type productQueryUsecase struct {
	repo        domain.ProductRepository
	defaultLang string
}

// NewProductQueryUsecase creates a new instance of product query business logic.
func NewProductQueryUsecase(repo domain.ProductRepository, defaultLang string) domain.ProductQueryUsecase {
	return &productQueryUsecase{
		repo:        repo,
		defaultLang: defaultLang,
	}
}

func (u *productQueryUsecase) GetAllProducts(ctx context.Context, p domain.Pagination) ([]domain.Product, error) {
	products, err := u.repo.FindAll(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetAllProducts: %w", err)
	}
	return products, nil
}

func (u *productQueryUsecase) GetProductByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Product, error) {
	product, err := u.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetProductByPublicID: %w", err)
	}
	if product == nil {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

func (u *productQueryUsecase) GetProductDetails(ctx context.Context, query domain.GetProductDetailsQuery) (*domain.ProductDetailsResponse, error) {
	langCode := query.LangCode
	if langCode == "" {
		langCode = u.defaultLang
	}

	product, err := u.repo.GetProductDetails(ctx, query.PublicID, langCode)
	if err != nil {
		return nil, fmt.Errorf("productQueryUsecase.GetProductDetails: %w", err)
	}
	if product == nil {
		return nil, domain.ErrProductNotFound
	}

	// Mapping to DTO
	resp := &domain.ProductDetailsResponse{
		PublicID:  product.PublicID,
		BasePrice: 0, // In a real app, this might come from a pricing service or variant
		Status:    string(product.Status),
	}

	if product.Translation != nil {
		resp.Name = product.Translation.Name
		resp.Description = product.Translation.Description
		resp.ShortDesc = product.Translation.ShortDesc
	} else {
		// Fallback to base product fields if translation is missing
		resp.Name = product.Name
		resp.Description = product.Description
		resp.ShortDesc = product.ShortDesc
	}

	if product.SEO != nil {
		resp.MetaTitle = product.SEO.MetaTitle
		resp.MetaDescription = product.SEO.MetaDescription
	}

	if len(product.Categories) > 0 {
		resp.Categories = make([]string, len(product.Categories))
		for i, c := range product.Categories {
			resp.Categories[i] = c.Name
		}
	}

	if len(product.Tags) > 0 {
		resp.Tags = make([]string, len(product.Tags))
		for i, t := range product.Tags {
			resp.Tags[i] = t.Name
		}
	}

	return resp, nil
}
