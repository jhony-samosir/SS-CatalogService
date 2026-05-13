package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type ImportJob struct {
	ID          int
	PublicID    uuid.UUID
	JobType     string
	Status      JobStatus
	FileURL     string
	ErrorLog    string
	TotalRows   int
	Processed   int
	CreatedBy   string
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// --- Interfaces ---

type ImportRepository interface {
	Create(ctx context.Context, job *ImportJob) error
	GetByPublicID(ctx context.Context, publicID uuid.UUID) (*ImportJob, error)
	UpdateStatus(ctx context.Context, id int, status JobStatus, errorLog string, processed int) error
	GetPendingJobs(ctx context.Context) ([]ImportJob, error)
}

type ImportUsecase interface {
	TriggerImport(ctx context.Context, fileURL string, jobType string, userID string) (*ImportJob, error)
	GetJobStatus(ctx context.Context, publicID uuid.UUID) (*ImportJob, error)
}
