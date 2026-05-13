package postgres

import (
	"context"
	"ss-catalog-service/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type importRepository struct {
	db *gorm.DB
}

func NewImportRepository(db *gorm.DB) domain.ImportRepository {
	return &importRepository{db: db}
}

func (r *importRepository) Create(ctx context.Context, job *domain.ImportJob) error {
	model := &ImportJobModel{
		JobType:   job.JobType,
		Status:    string(job.Status),
		FileURL:   job.FileURL,
		CreatedBy: job.CreatedBy,
	}
	db := getDB(ctx, r.db)
	if err := db.Create(model).Error; err != nil {
		return err
	}
	job.ID = model.ID
	job.PublicID = model.PublicID
	job.CreatedAt = model.CreatedAt
	return nil
}

func (r *importRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.ImportJob, error) {
	var model ImportJobModel
	db := getDB(ctx, r.db)
	if err := db.Where("public_id = ?", publicID).First(&model).Error; err != nil {
		return nil, err
	}
	job := model.ToDomain()
	return &job, nil
}

func (r *importRepository) UpdateStatus(ctx context.Context, id int, status domain.JobStatus, errorLog string, processed int) error {
	db := getDB(ctx, r.db)
	updates := map[string]interface{}{
		"status":    string(status),
		"error_log": errorLog,
		"processed": processed,
	}
	if status == domain.JobStatusCompleted || status == domain.JobStatusFailed {
		now := time.Now()
		updates["completed_at"] = &now
	}
	return db.Model(&ImportJobModel{}).Where("id = ?", id).Updates(updates).Error
}

func (r *importRepository) GetPendingJobs(ctx context.Context) ([]domain.ImportJob, error) {
	var models []ImportJobModel
	db := getDB(ctx, r.db)
	if err := db.Where("status = ?", string(domain.JobStatusPending)).Find(&models).Error; err != nil {
		return nil, err
	}
	jobs := make([]domain.ImportJob, len(models))
	for i, m := range models {
		jobs[i] = m.ToDomain()
	}
	return jobs, nil
}
