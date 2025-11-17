package repository

import (
	"context"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SessionRepository defines the data access methods for a session.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	ListSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Session, error)
	UpdateSessionStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status string) (*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
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
		INSERT INTO sessions (id, user_id, notebook_id, current_kernel_id, status, last_active_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, notebook_id, current_kernel_id, status, last_active_at;
	`
	row := r.db.QueryRow(ctx, query,
		session.ID,
		session.UserID,
		session.NotebookID,
		session.CurrentKernelID,
		session.Status,
		session.LastActiveAt,
	)

	var createdSession models.Session
	if err := row.Scan(
		&createdSession.ID,
		&createdSession.UserID,
		&createdSession.NotebookID,
		&createdSession.CurrentKernelID,
		&createdSession.Status,
		&createdSession.LastActiveAt,
	); err != nil {
		return nil, err
	}

	return &createdSession, nil
}

// ListSessions retrieves all sessions for a given user ID.
func (r *sessionRepository) ListSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	query := `
		SELECT id, user_id, notebook_id, current_kernel_id, status, last_active_at
		FROM sessions
		WHERE user_id = $1;
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.NotebookID,
			&session.CurrentKernelID,
			&session.Status,
			&session.LastActiveAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

// GetSessionByID retrieves a single session by its ID and user ID.
func (r *sessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.Session, error) {
	query := `
		SELECT id, user_id, notebook_id, current_kernel_id, status, last_active_at
		FROM sessions
		WHERE id = $1 AND user_id = $2;
	`
	row := r.db.QueryRow(ctx, query, id, userID)

	var session models.Session
	if err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.NotebookID,
		&session.CurrentKernelID,
		&session.Status,
		&session.LastActiveAt,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

// UpdateSessionStatus updates the status and last_active_at fields of a session.
func (r *sessionRepository) UpdateSessionStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status string) (*models.Session, error) {
	query := `
		UPDATE sessions
		SET status = $3, last_active_at = $4
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, notebook_id, current_kernel_id, status, last_active_at;
	`
	row := r.db.QueryRow(ctx, query, id, userID, status, time.Now().UTC())

	var updatedSession models.Session
	if err := row.Scan(
		&updatedSession.ID,
		&updatedSession.UserID,
		&updatedSession.NotebookID,
		&updatedSession.CurrentKernelID,
		&updatedSession.Status,
		&updatedSession.LastActiveAt,
	); err != nil {
		return nil, err
	}

	return &updatedSession, nil
}

// DeleteSession deletes a session from the database.
func (r *sessionRepository) DeleteSession(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1 AND user_id = $2;
	`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

