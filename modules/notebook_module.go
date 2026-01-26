package modules

import (
	"context"
	"errors"
	"fmt"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
)

// NotebookModule encapsulates business logic for notebooks.
type NotebookModule struct {
	repo        repository.NotebookRepository
	ProblemRepo repository.ProblemRepository // Added ProblemRepository
}

// NewNotebookModule creates and returns a new NotebookModule.
func NewNotebookModule(
	repo repository.NotebookRepository,
	problemRepo repository.ProblemRepository,
) *NotebookModule {
	return &NotebookModule{
		repo:        repo,
		ProblemRepo: problemRepo,
	}
}

// CreateNotebook handles the business logic for creating a new notebook.
func (m *NotebookModule) CreateNotebook(
	ctx context.Context,
	req *models.CreateNotebookRequest,
	userID string,
) (*models.Notebook, error) {
	if req == nil {
		return nil, errors.New("invalid notebook request")
	}

	// Verify ownership of the problem statement if provided
	if req.ProblemStatementID == nil || *req.ProblemStatementID == "" {
		return nil, errors.New("problem statement ID is required to create a notebook")
	}

	problem, err := m.ProblemRepo.GetProblemByID(ctx, *req.ProblemStatementID)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem statement for ownership verification: %w", err)
	}
	if problem == nil || problem.CreatedBy.String() != userID {
		return nil, errors.New("problem statement not found or not owned by user")
	}

	return m.repo.CreateNotebook(ctx, req)
}

// ListNotebooks handles the business logic for listing notebooks.
func (m *NotebookModule) ListNotebooks(
	ctx context.Context,
	filters map[string]string,
	userID string,
) ([]models.Notebook, error) {
	// Add userID to filters to ensure only owner's notebooks are listed
	if filters == nil {
		filters = make(map[string]string)
	}
	filters["created_by"] = userID // This filter will be handled in the repository

	return m.repo.ListNotebooks(ctx, filters, userID)
}

// GetNotebookByID handles the business logic for retrieving a single notebook.
func (m *NotebookModule) GetNotebookByID(
	ctx context.Context,
	id string,
	userID string,
) (*models.Notebook, error) {
	// The repository will enforce ownership.
	nb, err := m.repo.GetNotebookByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if nb == nil || nb.ID == "" { // Check if notebook was actually found
		return nil, errors.New("notebook not found or not owned by user")
	}
	return nb, nil
}

// UpdateNotebook handles the business logic for updating a notebook.
func (m *NotebookModule) UpdateNotebook(
	ctx context.Context,
	id string,
	req *models.UpdateNotebookRequest,
	userID string,
) (*models.Notebook, error) {
	if req == nil {
		return nil, errors.New("invalid update request")
	}

	// The repository will enforce ownership.
	updated, err := m.repo.UpdateNotebook(ctx, id, req, userID)
	if err != nil {
		return nil, err
	}
	if updated == nil || updated.ID == "" {
		return nil, errors.New("notebook not found or not owned by user")
	}

	return updated, nil
}

// DeleteNotebook handles the business logic for deleting a notebook.
func (m *NotebookModule) DeleteNotebook(
	ctx context.Context,
	id string,
	userID string,
) error {
	// The repository will enforce ownership.
	err := m.repo.DeleteNotebook(ctx, id, userID)
	if err != nil {
		return err
	}
	return nil
}
