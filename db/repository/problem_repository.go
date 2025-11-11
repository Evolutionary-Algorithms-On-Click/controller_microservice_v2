package repository

import (
	"context"
	"encoding/json"

	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ProblemRepository defines the interface for database operations on problem statements.
type ProblemRepository interface {
	CreateProblem(ctx context.Context, problem *models.ProblemStatement) (*models.ProblemStatement, error)
	GetProblemByID(ctx context.Context, problemID string) (*models.ProblemStatement, error)
	GetProblemsByUserID(ctx context.Context, userID string) ([]models.ProblemStatement, error)
	UpdateProblem(ctx context.Context, problemID string, title string, description json.RawMessage) (*models.ProblemStatement, error)
	DeleteProblem(ctx context.Context, problemID string) error
}

// problemRepository is the concrete implementation of ProblemRepository.
type problemRepository struct {
	db *pgxpool.Pool
}

// NewProblemRepository creates a new ProblemRepository.
func NewProblemRepository(db *pgxpool.Pool) ProblemRepository {
	return &problemRepository{db: db}
}

// CreateProblem inserts a new problem statement into the database.
func (r *problemRepository) CreateProblem(ctx context.Context, problem *models.ProblemStatement) (*models.ProblemStatement, error) {
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
		return nil, err
	}

	return &createdProblem, nil
}

// GetProblemByID retrieves a problem statement from the database by its ID.
func (r *problemRepository) GetProblemByID(ctx context.Context, problemID string) (*models.ProblemStatement, error) {
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
		return nil, err
	}

	return &problem, nil
}

// GetProblemsByUserID retrieves all problem statements for a given user.
func (r *problemRepository) GetProblemsByUserID(ctx context.Context, userID string) ([]models.ProblemStatement, error) {
	query := `
		SELECT id, title, description_json, created_by, created_at
		FROM problem_statements
		WHERE created_by = $1;
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
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
			return nil, err
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

// UpdateProblem updates a problem statement in the database.
func (r *problemRepository) UpdateProblem(ctx context.Context, problemID string, title string, description json.RawMessage) (*models.ProblemStatement, error) {
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
		return nil, err
	}

	return &updatedProblem, nil
}

// DeleteProblem deletes a problem statement from the database.
func (r *problemRepository) DeleteProblem(ctx context.Context, problemID string) error {
	query := `DELETE FROM problem_statements WHERE id = $1;`
	_, err := r.db.Exec(ctx, query, problemID)
	return err
}
