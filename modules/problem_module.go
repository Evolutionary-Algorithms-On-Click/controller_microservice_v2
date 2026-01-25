package modules

import (
	"context"
	"errors"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// ProblemModule encapsulates the business logic for problem statements.
type ProblemModule struct {
	ProblemRepo repository.ProblemRepository
	logger      zerolog.Logger
}

// NewProblemModule creates and returns a new ProblemModule.
func NewProblemModule(problemRepo repository.ProblemRepository, logger zerolog.Logger) *ProblemModule {
	return &ProblemModule{ProblemRepo: problemRepo, logger: logger}
}

// WithLogger allows setting the logger for the module.
func (m *ProblemModule) WithLogger(logger zerolog.Logger) *ProblemModule {
	m.logger = logger
	m.ProblemRepo = m.ProblemRepo.WithLogger(logger) // Also set logger for the repository
	return m
}

// CreateProblem inserts a new problem statement into the database.
func (m *ProblemModule) CreateProblem(ctx context.Context, problem *models.CreateProblemRequest, createdBy string) (*models.ProblemStatement, error) {
	if problem == nil {
		m.logger.Error().Msg("invalid problem creation request: problem is nil")
		return nil, errors.New("invalid problem creation request")
	}

	m.logger.Info().
		Str("title", problem.Title).
		Str("createdBy", createdBy).
		Msg("attempting to create new problem in module")

	creatorID, err := uuid.Parse(createdBy)
	if err != nil {
		m.logger.Error().Err(err).Str("createdBy", createdBy).Msg("invalid creator ID provided")
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
		m.logger.Error().Err(err).
			Str("problemID", newProblem.ID.String()).
			Str("title", newProblem.Title).
			Msg("failed to create problem in repository")
		return nil, err
	}

	m.logger.Info().
		Str("problemID", createdProblem.ID.String()).
		Msg("successfully created new problem in module")

	return createdProblem, nil
}

// GetProblemByID retrieves a problem statement from the database.
func (m *ProblemModule) GetProblemByID(ctx context.Context, problemID string) (*models.ProblemStatement, error) {
	m.logger.Info().Str("problemID", problemID).Msg("attempting to get problem by ID in module")

	problem, err := m.ProblemRepo.GetProblemByID(ctx, problemID)
	if err != nil {
		m.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to get problem by ID from repository")
		return nil, err
	}
	if problem == nil { // This case is handled in the repository, where it returns nil, nil for ErrNoRows
		m.logger.Warn().Str("problemID", problemID).Msg("problem not found in module")
		return nil, errors.New("problem not found")
	}

	m.logger.Info().Str("problemID", problemID).Msg("successfully retrieved problem by ID in module")
	return problem, nil
}

// GetProblemsByUserID retrieves all problem statements created by a specific user.
func (m *ProblemModule) GetProblemsByUserID(ctx context.Context, userID string) ([]models.ProblemStatement, error) {
	m.logger.Info().Str("userID", userID).Msg("attempting to get problems by user ID in module")

	problems, err := m.ProblemRepo.GetProblemsByUserID(ctx, userID)
	if err != nil {
		m.logger.Error().Err(err).Str("userID", userID).Msg("failed to get problems by user ID from repository")
		return nil, err
	}

	m.logger.Info().
		Str("userID", userID).
		Int("problem_count", len(problems)).
		Msg("successfully retrieved problems by user ID in module")

	return problems, nil
}

// DeleteProblem deletes a problem statement.
func (m *ProblemModule) DeleteProblem(ctx context.Context, problemID string, userID string) error {
	m.logger.Info().
		Str("problemID", problemID).
		Str("userID", userID).
		Msg("attempting to delete problem in module")

	// First, get the problem to ensure it exists and to check ownership.
	problem, err := m.ProblemRepo.GetProblemByID(ctx, problemID)
	if err != nil {
		m.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to get problem by ID for deletion check")
		return err // Could be not found, or other DB error.
	}
	if problem == nil { // This means problem was not found in repo
		m.logger.Warn().Str("problemID", problemID).Msg("problem not found for deletion")
		return errors.New("problem not found")
	}

	// Authorization: Check if the user trying to delete is the one who created it.
	if problem.CreatedBy.String() != userID {
		m.logger.Warn().
			Str("problemID", problemID).
			Str("requestingUserID", userID).
			Str("problemOwnerID", problem.CreatedBy.String()).
			Msg("user not authorized to delete this problem")
		return errors.New("user not authorized to delete this problem")
	}

	err = m.ProblemRepo.DeleteProblem(ctx, problemID)
	if err != nil {
		m.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to delete problem from repository")
		return err
	}

	m.logger.Info().Str("problemID", problemID).Msg("successfully deleted problem in module")
	return nil
}

// UpdateProblem updates a problem statement.
func (m *ProblemModule) UpdateProblem(ctx context.Context, problemID string, req *models.UpdateProblemRequest, userID string) (*models.ProblemStatement, error) {
	m.logger.Info().
		Str("problemID", problemID).
		Str("userID", userID).
		Msg("attempting to update problem in module")

	// First, get the problem to ensure it exists and to check ownership.
	problem, err := m.ProblemRepo.GetProblemByID(ctx, problemID)
	if err != nil {
		m.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to get problem by ID for update check")
		return nil, err // Could be not found, or other DB error.
	}
	if problem == nil { // This means problem was not found in repo
		m.logger.Warn().Str("problemID", problemID).Msg("problem not found for update")
		return nil, errors.New("problem not found")
	}

	// Authorization: Check if the user trying to update is the one who created it.
	if problem.CreatedBy.String() != userID {
		m.logger.Warn().
			Str("problemID", problemID).
			Str("requestingUserID", userID).
			Str("problemOwnerID", problem.CreatedBy.String()).
			Msg("user not authorized to update this problem")
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

	updatedProblem, err := m.ProblemRepo.UpdateProblem(ctx, problemID, title, description)
	if err != nil {
		m.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to update problem in repository")
		return nil, err
	}

	m.logger.Info().Str("problemID", updatedProblem.ID.String()).Msg("successfully updated problem in module")
	return updatedProblem, nil
}

