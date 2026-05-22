package import_job

import (
	"context"
	"log/slog"
	"ss-catalog-service/internal/domain"

	"github.com/google/uuid"
)

type importUsecase struct {
	repo   domain.ImportRepository
	logger *slog.Logger
}

func NewImportUsecase(repo domain.ImportRepository, logger *slog.Logger) domain.ImportUsecase {
	return &importUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *importUsecase) TriggerImport(ctx context.Context, fileURL string, jobType string, userID string) (*domain.ImportJob, error) {
	job := &domain.ImportJob{
		FileURL:   fileURL,
		JobType:   jobType,
		Status:    domain.JobStatusPending,
		CreatedBy: userID,
	}

	u.logger.InfoContext(ctx, "triggering import job", "file_url", fileURL, "job_type", jobType, "user_id", userID)

	if err := u.repo.Create(ctx, job); err != nil {
		u.logger.ErrorContext(ctx, "failed to create import job", "error", err, "file_url", fileURL, "job_type", jobType)
		return nil, err
	}

	u.logger.InfoContext(ctx, "import job triggered successfully", "job_id", job.PublicID, "file_url", fileURL, "job_type", jobType)
	return job, nil
}

func (u *importUsecase) GetJobStatus(ctx context.Context, publicID uuid.UUID) (*domain.ImportJob, error) {
	return u.repo.GetByPublicID(ctx, publicID)
}

func (u *importUsecase) GetAllJobs(ctx context.Context, p domain.Pagination) ([]domain.ImportJob, int64, error) {
	return u.repo.FindAll(ctx, p)
}
