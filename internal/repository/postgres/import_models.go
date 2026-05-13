package postgres

import (
	"ss-catalog-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ImportJobModel struct {
	ID          int       `gorm:"primaryKey"`
	PublicID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
	JobType     string    `gorm:"type:varchar(50);not null"`
	Status      string    `gorm:"type:varchar(20);default:'pending'"`
	FileURL     string    `gorm:"type:varchar(500)"`
	ErrorLog    string    `gorm:"type:text"`
	TotalRows   int       `gorm:"default:0"`
	Processed   int       `gorm:"default:0"`
	CreatedBy   string    `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CompletedAt *time.Time
}

func (ImportJobModel) TableName() string {
	return "import_jobs"
}

func (m *ImportJobModel) ToDomain() domain.ImportJob {
	return domain.ImportJob{
		ID:          m.ID,
		PublicID:    m.PublicID,
		JobType:     m.JobType,
		Status:      domain.JobStatus(m.Status),
		FileURL:     m.FileURL,
		ErrorLog:    m.ErrorLog,
		TotalRows:   m.TotalRows,
		Processed:   m.Processed,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		CompletedAt: m.CompletedAt,
	}
}
