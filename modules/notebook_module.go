package modules

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
)

// NotebookModule encapsulates business logic for notebooks.
type NotebookModule struct{}

// NewNotebookModule creates and returns a new NotebookModule.
func NewNotebookModule() *NotebookModule {
	return &NotebookModule{}
}

// CreateNotebook inserts a new notebook record into the database.
func (m *NotebookModule) CreateNotebook(ctx context.Context, req *models.CreateNotebookRequest) (*models.Notebook, error) {
	if req == nil {
		return nil, errors.New("invalid notebook request")
	}

	id := uuid.New().String()
	now := time.Now().UTC()

	query := `
		INSERT INTO notebooks (id, title, context_minio_url, problem_statement_id, created_at, last_modified_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, context_minio_url, problem_statement_id, created_at, last_modified_at;
	`

	row := db.Pool.QueryRow(ctx, query,
		id,
		req.Title,
		nil, // TODO: Should include logic for context minIO url
		req.ProblemStatementID,
		now,
		now,
	)

	var nb models.Notebook
	if err := row.Scan(
		&nb.ID,
		&nb.Title,
		&nb.ContextMinioURL,
		&nb.ProblemStatementID,
		&nb.CreatedAt,
		&nb.LastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &nb, nil
}

// ListNotebooks retrieves all notebooks (optionally filtered by problem_id, user_id, etc.)
func (m *NotebookModule) ListNotebooks(ctx context.Context, filters map[string]string) ([]models.Notebook, error) {
	query := `SELECT id, title, context_minio_url, problem_statement_id, created_at, last_modified_at FROM notebooks`
	args := []any{}
	where := ""

	// optional filters (future extension)
	if filters != nil {
		i := 1
		for key, val := range filters {
			if i == 1 {
				where += " WHERE "
			} else {
				where += " AND "
			}
			where += key + " = $" + strconv.Itoa(i)
			args = append(args, val)
			i++
		}
	}

	rows, err := db.Pool.Query(ctx, query+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notebooks []models.Notebook
	for rows.Next() {
		var nb models.Notebook
		if err := rows.Scan(
			&nb.ID,
			&nb.Title,
			&nb.ContextMinioURL,
			&nb.ProblemStatementID,
			&nb.CreatedAt,
			&nb.LastModifiedAt,
		); err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}

	return notebooks, nil
}

// GetNotebookByID retrieves a single notebook by ID.
func (m *NotebookModule) GetNotebookByID(ctx context.Context, id string) (*models.Notebook, error) {
	query := `
		SELECT id, title, context_minio_url, problem_statement_id, created_at, last_modified_at
		FROM notebooks
		WHERE id = $1;
	`
	row := db.Pool.QueryRow(ctx, query, id)

	var nb models.Notebook
	if err := row.Scan(
		&nb.ID,
		&nb.Title,
		&nb.ContextMinioURL,
		&nb.ProblemStatementID,
		&nb.CreatedAt,
		&nb.LastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &nb, nil
}

// UpdateNotebook updates a notebook’s title or context_minio_url.
func (m *NotebookModule) UpdateNotebook(ctx context.Context, id string, req *models.UpdateNotebookRequest) (*models.Notebook, error) {
	if req == nil {
		return nil, errors.New("invalid update request")
	}

	// Build dynamic update query
	setClause := ""
	args := []any{}
	argIndex := 1

	if req.Title != nil {
		setClause += "title = $" + strconv.Itoa(argIndex)
		args = append(args, *req.Title)
		argIndex++
	}

	// TODO: MinIO url should be got from other business logic
	// if req.ContextMinioURL != nil {
	// 	if len(setClause) > 0 {
	// 		setClause += ", "
	// 	}
	// 	setClause += "context_minio_url = $" + strconv.Itoa(argIndex)
	// 	args = append(args, *req.ContextMinioURL)
	// 	argIndex++
	// }

	if setClause == "" {
		return m.GetNotebookByID(ctx, id)
	}

	// Always update last_modified_at
	setClause += ", last_modified_at = $" + strconv.Itoa(argIndex)
	args = append(args, time.Now().UTC())
	argIndex++

	args = append(args, id)

	query := `
		UPDATE notebooks
		SET ` + setClause + `
		WHERE id = $` + strconv.Itoa(argIndex) + `
		RETURNING id, title, context_minio_url, problem_statement_id, created_at, last_modified_at;
	`

	row := db.Pool.QueryRow(ctx, query, args...)

	var nb models.Notebook
	if err := row.Scan(
		&nb.ID,
		&nb.Title,
		&nb.ContextMinioURL,
		&nb.ProblemStatementID,
		&nb.CreatedAt,
		&nb.LastModifiedAt,
	); err != nil {
		return nil, err
	}

	return &nb, nil
}

// DeleteNotebook removes a notebook by ID.
func (m *NotebookModule) DeleteNotebook(ctx context.Context, id string) error {
	cmd, err := db.Pool.Exec(ctx, `DELETE FROM notebooks WHERE id = $1;`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("notebook not found")
	}
	return nil
}
