package attribute

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type attributeUsecase struct {
	attrRepo domain.AttributeRepository
}

func NewAttributeUsecase(attrRepo domain.AttributeRepository) domain.AttributeUsecase {
	return &attributeUsecase{attrRepo: attrRepo}
}

func (u *attributeUsecase) GetAttributes(ctx context.Context, p domain.Pagination) ([]domain.ProductAttribute, error) {
	return u.attrRepo.FindAll(ctx, p)
}

func (u *attributeUsecase) GetAttributeByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.ProductAttribute, error) {
	return u.attrRepo.FindByPublicID(ctx, publicID)
}

func (u *attributeUsecase) CreateAttribute(ctx context.Context, attr *domain.ProductAttribute) error {
	if attr.PublicID == uuid.Nil {
		attr.PublicID = uuid.New()
	}
	return u.attrRepo.Create(ctx, attr)
}

func (u *attributeUsecase) UpdateAttribute(ctx context.Context, attr *domain.ProductAttribute) error {
	return u.attrRepo.Update(ctx, attr)
}

func (u *attributeUsecase) DeleteAttribute(ctx context.Context, publicID uuid.UUID) error {
	return u.attrRepo.Delete(ctx, publicID)
}

// Tag Usecase Implementation
type tagUsecase struct {
	repo domain.TagRepository
}

func NewTagUsecase(repo domain.TagRepository) domain.TagUsecase {
	return &tagUsecase{repo: repo}
}

func (u *tagUsecase) GetTags(ctx context.Context, p domain.Pagination) ([]domain.Tag, error) {
	return u.repo.FindAll(ctx, p)
}

func (u *tagUsecase) CreateTag(ctx context.Context, tag *domain.Tag) error {
	if tag.PublicID == uuid.Nil {
		tag.PublicID = uuid.New()
	}
	return u.repo.Create(ctx, tag)
}

func (u *tagUsecase) DeleteTag(ctx context.Context, publicID uuid.UUID) error {
	return u.repo.Delete(ctx, publicID)
}
