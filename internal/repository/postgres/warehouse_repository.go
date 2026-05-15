package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) domain.WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) FindAll(ctx context.Context, p domain.Pagination) ([]domain.Warehouse, error) {
	var models []WarehouseModel
	db := getDB(ctx, r.db)

	query := db.Model(&WarehouseModel{})
	if p.Limit > 0 {
		query = query.Limit(p.Limit).Offset(p.Offset)
	}

	if err := query.Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	whs := make([]domain.Warehouse, len(models))
	for i, m := range models {
		whs[i] = domain.Warehouse{
			BaseEntity: domain.BaseEntity{
				ID:        m.ID,
				PublicID:  m.PublicID,
				CreatedAt: m.CreatedAt,
			},
			SellerID:    m.SellerID,
			Name:        m.Name,
			Code:        m.Code,
			City:        m.City,
			Province:    m.Province,
			CountryCode: m.CountryCode,
			PostalCode:  m.PostalCode,
			Address:     m.Address,
			IsActive:    m.IsActive,
		}
	}
	return whs, nil
}

func (r *warehouseRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.Warehouse, error) {
	var m WarehouseModel
	db := getDB(ctx, r.db)

	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Warehouse{
		BaseEntity: domain.BaseEntity{
			ID:        m.ID,
			PublicID:  m.PublicID,
		},
		SellerID:    m.SellerID,
		Name:        m.Name,
		Code:        m.Code,
		City:        m.City,
		Province:    m.Province,
		CountryCode: m.CountryCode,
		PostalCode:  m.PostalCode,
		Address:     m.Address,
		IsActive:    m.IsActive,
	}, nil
}

func (r *warehouseRepository) Create(ctx context.Context, wh *domain.Warehouse) error {
	model := &WarehouseModel{
		BaseModel: BaseModel{PublicID: wh.PublicID},
		SellerID:    wh.SellerID,
		Name:        wh.Name,
		Code:        wh.Code,
		City:        wh.City,
		Province:    wh.Province,
		CountryCode: wh.CountryCode,
		PostalCode:  wh.PostalCode,
		Address:     wh.Address,
		IsActive:    wh.IsActive,
	}
	db := getDB(ctx, r.db)
	return db.Create(model).Error
}

func (r *warehouseRepository) Update(ctx context.Context, wh *domain.Warehouse) error {
	db := getDB(ctx, r.db)
	user, _ := domain.UserFromContext(ctx)
	return db.Model(&WarehouseModel{}).
		Where("public_id = ?", wh.PublicID).
		Updates(map[string]interface{}{
			"seller_id":    wh.SellerID,
			"name":         wh.Name,
			"code":         wh.Code,
			"city":         wh.City,
			"province":     wh.Province,
			"country_code": wh.CountryCode,
			"postal_code":  wh.PostalCode,
			"address":      wh.Address,
			"is_active":    wh.IsActive,
			"updated_by":   user.FullName,
		}).Error
}

func (r *warehouseRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	db := getDB(ctx, r.db)
	var m WarehouseModel
	if err := db.Where("public_id = ?", publicID).First(&m).Error; err != nil {
		return err
	}
	return db.Delete(&m).Error
}

func (r *warehouseRepository) CountInventory(ctx context.Context, warehouseID int) (int64, error) {
	var count int64
	db := getDB(ctx, r.db)
	if err := db.Model(&ProductInventoryModel{}).Where("warehouse_id = ?", warehouseID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *warehouseRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	db := getDB(ctx, r.db)
	if err := db.Model(&WarehouseModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
