package repository

import (
	"context"
	"encoding/json"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
)

// ProblemRepository defines the interface for database operations on problem statements.
type ProblemRepository interface {
	CreateProblem(ctx context.Context, problem *models.ProblemStatement) (*models.ProblemStatement, error)
	GetProblemByID(ctx context.Context, problemID string) (*models.ProblemStatement, error)
	GetProblemsByUserID(ctx context.Context, userID string) ([]models.ProblemStatement, error)
	UpdateProblem(ctx context.Context, problemID string, title string, description json.RawMessage) (*models.ProblemStatement, error)
	DeleteProblem(ctx context.Context, problemID string) error
	WithLogger(logger zerolog.Logger) ProblemRepository // Add WithLogger method
}

// problemRepository is the concrete implementation of ProblemRepository.
type problemRepository struct {
	db     *pgxpool.Pool
	logger zerolog.Logger // Add logger field
}

// NewProblemRepository creates a new ProblemRepository.
func NewProblemRepository(db *pgxpool.Pool) ProblemRepository {
	return &problemRepository{db: db}
}

// WithLogger allows setting the logger for the repository.
func (r *problemRepository) WithLogger(logger zerolog.Logger) ProblemRepository {
	r.logger = logger
	return r
}

// CreateProblem inserts a new problem statement into the database.
func (r *problemRepository) CreateProblem(ctx context.Context, problem *models.ProblemStatement) (*models.ProblemStatement, error) {
	r.logger.Info().
		Str("problemID", problem.ID.String()).
		Str("title", problem.Title).
		Msg("attempting to create new problem")

	query := `
		INSERT INTO problem_statements (id, title, description_json, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description_json, created_by, created_at;
	`
	row := r.db.QueryRow(ctx, query,
		problem.ID,
		problem.Title,
		problem.DescriptionJSON,
		problem.CreatedBy,
		problem.CreatedAt,
	)

	var createdProblem models.ProblemStatement
	if err := row.Scan(
		&createdProblem.ID,
		&createdProblem.Title,
		&createdProblem.DescriptionJSON,
		&createdProblem.CreatedBy,
		&createdProblem.CreatedAt,
	); err != nil {
		r.logger.Error().Err(err).
			Str("problemID", problem.ID.String()).
			Str("title", problem.Title).
			Msg("failed to create problem or scan created problem")
		return nil, err
	}

	r.logger.Info().
		Str("problemID", createdProblem.ID.String()).
		Msg("successfully created new problem")

	return &createdProblem, nil
}

// GetProblemByID retrieves a problem statement from the database by its ID.
func (r *problemRepository) GetProblemByID(ctx context.Context, problemID string) (*models.ProblemStatement, error) {
	r.logger.Info().Str("problemID", problemID).Msg("attempting to get problem by ID")

	query := `
		SELECT id, title, description_json, created_by, created_at
		FROM problem_statements
		WHERE id = $1;
	`
	row := r.db.QueryRow(ctx, query, problemID)

	var problem models.ProblemStatement
	if err := row.Scan(
		&problem.ID,
		&problem.Title,
		&problem.DescriptionJSON,
		&problem.CreatedBy,
		&problem.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn().Str("problemID", problemID).Msg("problem not found")
			return nil, nil // Or return a custom error like models.ErrNotFound
		}
		r.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to get problem by ID or scan row")
		return nil, err
	}

	r.logger.Info().Str("problemID", problemID).Msg("successfully retrieved problem by ID")
	return &problem, nil
}

// GetProblemsByUserID retrieves all problem statements for a given user.
func (r *problemRepository) GetProblemsByUserID(ctx context.Context, userID string) ([]models.ProblemStatement, error) {
	r.logger.Info().Str("userID", userID).Msg("attempting to get problems by user ID")

	query := `
		SELECT id, title, description_json, created_by, created_at
		FROM problem_statements
		WHERE created_by = $1;
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("failed to query problems by user ID")
		return nil, err
	}
	defer rows.Close()

	var problems []models.ProblemStatement
	for rows.Next() {
		var problem models.ProblemStatement
		if err := rows.Scan(
			&problem.ID,
			&problem.Title,
			&problem.DescriptionJSON,
			&problem.CreatedBy,
			&problem.CreatedAt,
		); err != nil {
			r.logger.Error().Err(err).Str("userID", userID).Msg("failed to scan problem row by user ID")
			return nil, err
		}
		problems = append(problems, problem)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("error after iterating through problem rows by user ID")
		return nil, err
	}

	r.logger.Info().
		Str("userID", userID).
		Int("problem_count", len(problems)).
		Msg("successfully retrieved problems by user ID")

	return problems, nil
}

// UpdateProblem updates a problem statement in the database.
func (r *problemRepository) UpdateProblem(ctx context.Context, problemID string, title string, description json.RawMessage) (*models.ProblemStatement, error) {
	r.logger.Info().
		Str("problemID", problemID).
		Str("title", title).
		Msg("attempting to update problem")

	query := `
		UPDATE problem_statements
		SET title = $2, description_json = $3
		WHERE id = $1
		RETURNING id, title, description_json, created_by, created_at;
	`
	row := r.db.QueryRow(ctx, query, problemID, title, description)

	var updatedProblem models.ProblemStatement
	if err := row.Scan(
		&updatedProblem.ID,
		&updatedProblem.Title,
		&updatedProblem.DescriptionJSON,
		&updatedProblem.CreatedBy,
		&updatedProblem.CreatedAt,
	); err != nil {
		r.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to update problem or scan updated problem")
		return nil, err
	}

	r.logger.Info().
		Str("problemID", updatedProblem.ID.String()).
		Msg("successfully updated problem")

	return &updatedProblem, nil
}

// DeleteProblem deletes a problem statement from the database.
func (r *problemRepository) DeleteProblem(ctx context.Context, problemID string) error {
	r.logger.Info().Str("problemID", problemID).Msg("attempting to delete problem")

	query := `DELETE FROM problem_statements WHERE id = $1;`
	_, err := r.db.Exec(ctx, query, problemID)
	if err != nil {
		r.logger.Error().Err(err).Str("problemID", problemID).Msg("failed to delete problem")
		return err
	}

	r.logger.Info().Str("problemID", problemID).Msg("successfully deleted problem")
	return nil
}
