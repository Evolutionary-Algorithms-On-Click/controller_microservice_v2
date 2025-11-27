package modules

import (
	"context"
	"errors"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
)

// NotebookModule encapsulates business logic for notebooks.
type NotebookModule struct {
	repo repository.NotebookRepository
}

// NewNotebookModule creates and returns a new NotebookModule.
func NewNotebookModule(repo repository.NotebookRepository) *NotebookModule {
	return &NotebookModule{
		repo: repo,
	}
}

// CreateNotebook handles the business logic for creating a new notebook.
func (m *NotebookModule) CreateNotebook(ctx context.Context, req *models.CreateNotebookRequest) (*models.Notebook, error) {
	if req == nil {
		return nil, errors.New("invalid notebook request")
	}
	// Any other business logic before creation would go here.
	return m.repo.CreateNotebook(ctx, req)
}

// ListNotebooks handles the business logic for listing notebooks.
func (m *NotebookModule) ListNotebooks(ctx context.Context, filters map[string]string) ([]models.Notebook, error) {
	// Business logic for filtering or pagination could go here.
	return m.repo.ListNotebooks(ctx, filters)
}

// GetNotebookByID handles the business logic for retrieving a single notebook.
func (m *NotebookModule) GetNotebookByID(ctx context.Context, id string) (*models.Notebook, error) {
	// Business logic for access control, etc., could go here.
	return m.repo.GetNotebookByID(ctx, id)
}

// UpdateNotebook handles the business logic for updating a notebook.
func (m *NotebookModule) UpdateNotebook(ctx context.Context, id string, req *models.UpdateNotebookRequest) (*models.Notebook, error) {
	if req == nil {
		return nil, errors.New("invalid update request")
	}

	// Business logic: if there's nothing to update, just return the current state.
	if req.Title == nil {
		return m.repo.GetNotebookByID(ctx, id)
	}

	return m.repo.UpdateNotebook(ctx, id, req)
}

// DeleteNotebook handles the business logic for deleting a notebook.
func (m *NotebookModule) DeleteNotebook(ctx context.Context, id string) error {
	// Business logic for checking permissions, etc., could go here.
	return m.repo.DeleteNotebook(ctx, id)
}