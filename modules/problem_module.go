package modules

import (
	"context"
	"errors"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
)

// ProblemModule encapsulates the business logic for problem statements.
type ProblemModule struct {
	ProblemRepo repository.ProblemRepository
}

// NewProblemModule creates and returns a new ProblemModule.
func NewProblemModule(problemRepo repository.ProblemRepository) *ProblemModule {
	return &ProblemModule{ProblemRepo: problemRepo}
}

// CreateProblem inserts a new problem statement into the database.
func (m *ProblemModule) CreateProblem(ctx context.Context, problem *models.CreateProblemRequest, createdBy string) (*models.ProblemStatement, error) {
	if problem == nil {
		return nil, errors.New("invalid problem creation request")
	}

	creatorID, err := uuid.Parse(createdBy)
	if err != nil {
		// If createdBy is an empty string or not a valid UUID, handle it.
		// For now, let's return an error. In a real app, you might have a default user or other logic.
		return nil, errors.New("invalid creator ID")
	}

	newProblem := &models.ProblemStatement{
		ID:              uuid.New(),
		Title:           problem.Title,
		DescriptionJSON: problem.DescriptionJSON,
		CreatedBy:       creatorID,
		CreatedAt:       time.Now().UTC(),
	}

	createdProblem, err := m.ProblemRepo.CreateProblem(ctx, newProblem)
	if err != nil {
		return nil, err
	}

	return createdProblem, nil
}

// GetProblemByID retrieves a problem statement from the database.
func (m *ProblemModule) GetProblemByID(ctx context.Context, problemID string) (*models.ProblemStatement, error) {
	return m.ProblemRepo.GetProblemByID(ctx, problemID)
}

// GetProblemsByUserID retrieves all problem statements created by a specific user.
func (m *ProblemModule) GetProblemsByUserID(ctx context.Context, userID string) ([]models.ProblemStatement, error) {
	return m.ProblemRepo.GetProblemsByUserID(ctx, userID)
}

// DeleteProblem deletes a problem statement.
func (m *ProblemModule) DeleteProblem(ctx context.Context, problemID string, userID string) error {
	// First, get the problem to ensure it exists and to check ownership.
	problem, err := m.ProblemRepo.GetProblemByID(ctx, problemID)
	if err != nil {
		return err // Could be not found, or other DB error.
	}

	// Authorization: Check if the user trying to delete is the one who created it.
	if problem.CreatedBy.String() != userID {
		return errors.New("user not authorized to delete this problem")
	}

	return m.ProblemRepo.DeleteProblem(ctx, problemID)
}

// UpdateProblem updates a problem statement.
func (m *ProblemModule) UpdateProblem(ctx context.Context, problemID string, req *models.UpdateProblemRequest, userID string) (*models.ProblemStatement, error) {
	// First, get the problem to ensure it exists and to check ownership.
	problem, err := m.ProblemRepo.GetProblemByID(ctx, problemID)
	if err != nil {
		return nil, err // Could be not found, or other DB error.
	}

	// Authorization: Check if the user trying to update is the one who created it.
	if problem.CreatedBy.String() != userID {
		return nil, errors.New("user not authorized to update this problem")
	}

	// If the request doesn't specify a new title, use the existing one.
	title := problem.Title
	if req.Title != "" {
		title = req.Title
	}

	// If the request doesn't specify a new description, use the existing one.
	description := problem.DescriptionJSON
	if req.DescriptionJSON != nil {
		description = req.DescriptionJSON
	}

	return m.ProblemRepo.UpdateProblem(ctx, problemID, title, description)
}

