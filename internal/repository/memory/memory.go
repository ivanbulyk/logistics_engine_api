package memory

import (
	"context"
	"sync"

	"github.com/ivanbulyk/logistics_engine_api/internal/model"
	"github.com/ivanbulyk/logistics_engine_api/internal/repository"
)

// Repository defines a memory allocation service repository
type Repository struct {
	DB sync.Map
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{}
}

// GetAll returns all report data.
func (r *Repository) GetAll(_ context.Context) ([]model.MetricsReport, error) {
	var reports []model.MetricsReport

	r.DB.Range(func(key, value interface{}) bool {
		reports = append(reports, value.(model.MetricsReport))
		return true
	})

	return reports, nil
}

// GetByID returns report data by id.
func (r *Repository) GetByID(_ context.Context, id int64) (model.MetricsReport, error) {

	if report, exist := r.DB.Load(id); exist {
		return report.(model.MetricsReport), nil
	}

	return model.MetricsReport{}, repository.ErrNotFound
}

// Create creates report .
func (r *Repository) Create(_ context.Context, report model.MetricsReport) (model.MetricsReport, error) {

	if _, exist := r.DB.Load(report.ID); exist {
		return model.MetricsReport{}, repository.ErrAlreadyExists
	}

	r.DB.Store(report.ID, report)

	return report, nil
}

// Update updates report data.
func (r *Repository) Update(_ context.Context, report model.MetricsReport) error {

	if _, exist := r.DB.Load(report.ID); !exist {
		return repository.ErrNotFound
	}

	r.DB.Store(report.ID, report)

	return nil
}

// Delete deletes report data.
func (r *Repository) Delete(_ context.Context, id int64) error {

	if _, exist := r.DB.Load(id); !exist {
		return repository.ErrNotFound
	}

	r.DB.Delete(id)

	return nil
}
