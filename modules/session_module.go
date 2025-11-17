package modules

import (
	"context"
	"errors"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	jupyterclient "github.com/Thanus-Kumaar/controller_microservice_v2/pkg/jupyter_client"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// SessionModule encapsulates the business logic for sessions.
type SessionModule struct {
	Repo    repository.SessionRepository
	Jupyter *jupyterclient.Client
	Logger  zerolog.Logger
}

// NewSessionModule creates and returns a new SessionModule.
func NewSessionModule(repo repository.SessionRepository, jupyter *jupyterclient.Client, logger zerolog.Logger) *SessionModule {
	return &SessionModule{
		Repo:    repo,
		Jupyter: jupyter,
		Logger:  logger,
	}
}

// CreateSession starts a new kernel and creates a session record in the database.
func (m *SessionModule) CreateSession(ctx context.Context, userIDStr string, notebookIDStr string, language string) (*models.Session, error) {
	if m.Jupyter == nil {
		return nil, errors.New("Jupyter client is not initialized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	notebookID, err := uuid.Parse(notebookIDStr)
	if err != nil {
		return nil, errors.New("invalid notebook_id format")
	}

	// TODO: Should check if the language is supported by the jupyter kernelspecs
	kernel, err := m.Jupyter.StartKernel(ctx, language)
	if err != nil {
		return nil, err
	}

	kernelID, err := uuid.Parse(kernel.ID)
	if err != nil {
		return nil, errors.New("invalid kernel_id format from Jupyter gateway")
	}

	newSession := &models.Session{
		ID:              uuid.New(),
		UserID:          userID, // Assign the parsed userID
		NotebookID:      notebookID,
		CurrentKernelID: kernelID,
		Status:          "active",
		LastActiveAt:    time.Now().UTC(),
	}

	createdSession, err := m.Repo.CreateSession(ctx, newSession)
	if err != nil {
		// If DB write fails, we should try to kill the orphaned kernel.
		m.Logger.Error().Err(err).Msg("failed to create session in DB, attempting to delete orphaned kernel")
		if deleteErr := m.Jupyter.DeleteKernel(context.Background(), kernel.ID); deleteErr != nil {
			m.Logger.Error().Err(deleteErr).Str("kernel_id", kernel.ID).Msg("failed to delete orphaned kernel")
		}
		return nil, err
	}

	return createdSession, nil
}

// ListSessions retrieves all sessions for a given user.
func (m *SessionModule) ListSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	sessions, err := m.Repo.ListSessions(ctx, userID)
	if err != nil {
		m.Logger.Error().Err(err).Msg("failed to list sessions from repo")
		return nil, err
	}

	return sessions, nil
}

// GetSessionByID retrieves a single session by its ID for a given user.
func (m *SessionModule) GetSessionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Session, error) {
	session, err := m.Repo.GetSessionByID(ctx, id, userID)
	if err != nil {
		m.Logger.Error().Err(err).Msg("failed to get session by id from repo")
		return nil, err
	}

	return session, nil
}

// UpdateSessionStatus updates the status of a session.
func (m *SessionModule) UpdateSessionStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status string) (*models.Session, error) {
	session, err := m.Repo.UpdateSessionStatus(ctx, id, userID, status)
	if err != nil {
		m.Logger.Error().Err(err).Msg("failed to update session status in repo")
		return nil, err
	}

	return session, nil
}

// DeleteSession deletes a session.
func (m *SessionModule) DeleteSession(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	err := m.Repo.DeleteSession(ctx, id, userID)
	if err != nil {
		m.Logger.Error().Err(err).Msg("failed to delete session from repo")
		return err
	}

	return nil
}
