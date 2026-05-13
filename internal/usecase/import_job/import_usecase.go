package import_job

import (
	"context"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type importUsecase struct {
	repo domain.ImportRepository
}

func NewImportUsecase(repo domain.ImportRepository) domain.ImportUsecase {
	return &importUsecase{repo: repo}
}

func (u *importUsecase) TriggerImport(ctx context.Context, fileURL string, jobType string, userID string) (*domain.ImportJob, error) {
	job := &domain.ImportJob{
		FileURL:   fileURL,
		JobType:   jobType,
		Status:    domain.JobStatusPending,
		CreatedBy: userID,
	}
	if err := u.repo.Create(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (u *importUsecase) GetJobStatus(ctx context.Context, publicID uuid.UUID) (*domain.ImportJob, error) {
	return u.repo.GetByPublicID(ctx, publicID)
}
