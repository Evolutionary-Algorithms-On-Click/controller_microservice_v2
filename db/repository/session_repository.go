package repository

import (
	"context"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SessionRepository defines the data access methods for a session.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
}

// sessionRepository is the concrete implementation of SessionRepository.
type sessionRepository struct {
	db *pgxpool.Pool
}

// NewSessionRepository creates a new SessionRepository.
func NewSessionRepository(db *pgxpool.Pool) SessionRepository {
	return &sessionRepository{db: db}
}

// CreateSession inserts a new session into the database.
func (r *sessionRepository) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	query := `
		INSERT INTO sessions (id, notebook_id, current_kernel_id, status, last_active_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, notebook_id, current_kernel_id, status, last_active_at;
	`
	row := r.db.QueryRow(ctx, query,
		session.ID,
		session.NotebookID,
		session.CurrentKernelID,
		session.Status,
		session.LastActiveAt,
	)

	var createdSession models.Session
	if err := row.Scan(
		&createdSession.ID,
		&createdSession.NotebookID,
		&createdSession.CurrentKernelID,
		&createdSession.Status,
		&createdSession.LastActiveAt,
	); err != nil {
		return nil, err
	}

	return &createdSession, nil
}

