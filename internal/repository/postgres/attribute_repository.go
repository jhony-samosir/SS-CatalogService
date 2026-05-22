package postgres

import (
	"context"
	"log/slog"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type attributeRepository struct {
	db *gorm.DB
}

func NewAttributeRepository(db *gorm.DB) domain.AttributeRepository {
	return &attributeRepository{db: db}
}

func (r *attributeRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.ProductAttribute, error) {
	var models []ProductAttributeModel
	db := getDB(ctx, r.db)

	query := db.Model(&ProductAttributeModel{}).Preload("Values")
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Order("sort_order ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	attrs := make([]domain.ProductAttribute, len(models))
	for i, m := range models {
		vals := make([]domain.AttributeValue, len(m.Values))
		for j, v := range m.Values {
			vals[j] = domain.AttributeValue{
				BaseEntity: domain.BaseEntity{
					ID:        v.ID,
					PublicID:  v.PublicID,
					CreatedAt: v.CreatedAt,
				},
				AttributeID: v.AttributeID,
				Value:       v.Value,
				ColorHex:    v.ColorHex,
				SortOrder:   v.SortOrder,
			}
		}

		attrs[i] = domain.ProductAttribute{
			BaseEntity: domain.BaseEntity{
				ID:        m.ID,
				PublicID:  m.PublicID,
				CreatedAt: m.CreatedAt,
			},
			Name:      m.Name,
			Code:      m.Code,
			InputType: domain.AttributeInputType(m.InputType),
			IsVariant: m.IsVariant,
			SortOrder: m.SortOrder,
			Values:    vals,
		}
	}
	return attrs, nil
}

func (r *attributeRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.ProductAttribute, error) {
	var m ProductAttributeModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).Preload("Values").First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	vals := make([]domain.AttributeValue, len(m.Values))
	for j, v := range m.Values {
		vals[j] = domain.AttributeValue{
			BaseEntity: domain.BaseEntity{ID: v.ID, PublicID: v.PublicID},
			AttributeID: v.AttributeID,
			Value:       v.Value,
			ColorHex:    v.ColorHex,
			SortOrder:   v.SortOrder,
		}
	}

	return &domain.ProductAttribute{
		BaseEntity: domain.BaseEntity{ID: m.ID, PublicID: m.PublicID},
		Name:      m.Name,
		Code:      m.Code,
		InputType: domain.AttributeInputType(m.InputType),
		IsVariant: m.IsVariant,
		SortOrder: m.SortOrder,
		Values:    vals,
	}, nil
}

func (r *attributeRepository) Create(ctx context.Context, attr *domain.ProductAttribute) error {
	model := &ProductAttributeModel{
		BaseModel: BaseModel{PublicID: attr.PublicID},
		Name:      attr.Name,
		Code:      attr.Code,
		InputType: string(attr.InputType),
		IsVariant: attr.IsVariant,
		SortOrder: attr.SortOrder,
	}
	db := getDB(ctx, r.db)
	err := db.Create(model).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_attribute", "error", err)
		return err
	}
	attr.ID = model.ID
	attr.CreatedAt = model.CreatedAt
	slog.InfoContext(ctx, "db_audit", "operation", "create_attribute", "attribute_id", attr.ID, "code", attr.Code)
	return nil
}

func (r *attributeRepository) Update(ctx context.Context, attr *domain.ProductAttribute) error {
	db := getDB(ctx, r.db)
	user, _ := domain.UserFromContext(ctx)
	err := db.Model(&ProductAttributeModel{}).
		Where("public_id = ?", attr.PublicID).
		Updates(map[string]interface{}{
			"name":       attr.Name,
			"code":       attr.Code,
			"input_type": string(attr.InputType),
			"is_variant": attr.IsVariant,
			"sort_order": attr.SortOrder,
			"updated_by": user.FullName,
		}).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "update_attribute", "public_id", attr.PublicID, "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "update_attribute", "public_id", attr.PublicID, "code", attr.Code)
	return nil
}

func (r *attributeRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	db := getDB(ctx, r.db)
	var m ProductAttributeModel
	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "delete_attribute_find", "public_id", publicID, "error", err)
		return err
	}
	err := db.Delete(&m).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "delete_attribute_execute", "public_id", publicID, "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "delete_attribute", "public_id", publicID, "attribute_id", m.ID)
	return nil
}

func (r *attributeRepository) CountUsage(ctx context.Context, attrID int) (int64, error) {
	var count int64
	db := getDB(ctx, r.db)
	if err := db.Model(&ProductVariantAttributeModel{}).Where("attribute_id = ?", attrID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *attributeRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	db := getDB(ctx, r.db)
	if err := db.Model(&ProductAttributeModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *attributeRepository) CreateValue(ctx context.Context, val *domain.AttributeValue) error {
	model := &AttributeValueModel{
		BaseModel:   BaseModel{PublicID: val.PublicID},
		AttributeID: val.AttributeID,
		Value:       val.Value,
		ColorHex:    val.ColorHex,
		SortOrder:   val.SortOrder,
	}
	db := getDB(ctx, r.db)
	err := db.Create(model).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_attribute_value", "attribute_id", val.AttributeID, "error", err)
		return err
	}
	val.ID = model.ID
	val.CreatedAt = model.CreatedAt
	slog.InfoContext(ctx, "db_audit", "operation", "create_attribute_value", "value_id", val.ID, "attribute_id", val.AttributeID, "value", val.Value)
	return nil
}

func (r *attributeRepository) DeleteValue(ctx context.Context, valID int) error {
	db := getDB(ctx, r.db)
	err := db.Delete(&AttributeValueModel{}, valID).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "delete_attribute_value", "value_id", valID, "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "delete_attribute_value", "value_id", valID)
	return nil
}

// Tag Repository Implementation

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) domain.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Tag, int64, error) {
	var models []TagModel
	var total int64
	db := getDB(ctx, r.db)
	query := db.Model(&TagModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}
	tags := make([]domain.Tag, len(models))
	for i, m := range models {
		tags[i] = domain.Tag{
			BaseEntity: domain.BaseEntity{ID: m.ID, PublicID: m.PublicID},
			Name:       m.Name,
			Slug:       m.Slug,
		}
	}
	return tags, total, nil
}

func (r *tagRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Tag, error) {
	var m TagModel
	db := getDB(ctx, r.db)
	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Tag{
		BaseEntity: domain.BaseEntity{ID: m.ID, PublicID: m.PublicID},
		Name:       m.Name,
		Slug:       m.Slug,
	}, nil
}

func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	model := &TagModel{
		BaseModel: BaseModel{PublicID: tag.PublicID},
		Name:      tag.Name,
		Slug:      tag.Slug,
	}
	db := getDB(ctx, r.db)
	err := db.Create(model).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "create_tag", "error", err)
		return err
	}
	tag.ID = model.ID
	slog.InfoContext(ctx, "db_audit", "operation", "create_tag", "tag_id", tag.ID, "slug", tag.Slug)
	return nil
}

func (r *tagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	db := getDB(ctx, r.db)
	user, _ := domain.UserFromContext(ctx)
	err := db.Model(&TagModel{}).
		Where("public_id = ?", tag.PublicID).
		Updates(map[string]interface{}{
			"name":       tag.Name,
			"slug":       tag.Slug,
			"updated_by": user.FullName,
		}).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "update_tag", "public_id", tag.PublicID, "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "update_tag", "public_id", tag.PublicID, "slug", tag.Slug)
	return nil
}

func (r *tagRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	db := getDB(ctx, r.db)
	var m TagModel
	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "delete_tag_find", "public_id", publicID, "error", err)
		return err
	}
	err := db.Delete(&m).Error
	if err != nil {
		slog.ErrorContext(ctx, "db_error", "operation", "delete_tag_execute", "public_id", publicID, "error", err)
		return err
	}
	slog.InfoContext(ctx, "db_audit", "operation", "delete_tag", "public_id", publicID, "tag_id", m.ID)
	return nil
}
