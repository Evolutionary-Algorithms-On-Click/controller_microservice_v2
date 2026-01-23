package modules

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Thanus-Kumaar/controller_microservice_v2/db/repository"
	"github.com/Thanus-Kumaar/controller_microservice_v2/pkg/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// CellModule encapsulates the business logic for cells.
type CellModule struct {
	Repo repository.CellRepository
	Logger zerolog.Logger
}

// NewCellModule creates and returns a new CellModule.
func NewCellModule(repo repository.CellRepository, logger zerolog.Logger) *CellModule {
	return &CellModule{
		Repo: repo,
		Logger: logger,
	}
}

func (m *CellModule) CreateCell(ctx context.Context, req *models.CreateCellRequest) (*models.Cell, error) {
	if req == nil {
		return nil, errors.New("invalid create cell request")
	}

	cell := &models.Cell{
		ID:         models.StringUUID(uuid.New()),
		NotebookID: req.NotebookID,
		CellIndex:  req.CellIndex,
		CellName: sql.NullString{
			String: req.CellName,
			Valid:  req.CellName != "",
		},
		CellType: req.CellType,
		Source:   req.Source,
	}

	return m.Repo.CreateCell(ctx, cell)
}

func (m *CellModule) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	return m.Repo.GetCellByID(ctx, id)
}

func (m *CellModule) GetCellsByNotebookID(ctx context.Context, notebookID uuid.UUID) ([]*models.Cell, error) {
	return m.Repo.GetCellsByNotebookID(ctx, notebookID)
}

func (m *CellModule) UpdateCell(ctx context.Context, id uuid.UUID, req *models.UpdateCellRequest) (*models.Cell, error) {
	cell, err := m.Repo.GetCellByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CellIndex != nil {
		cell.CellIndex = *req.CellIndex
	}
	if req.CellName != nil {
		cell.CellName.String = *req.CellName
		cell.CellName.Valid = true
	}
	if req.CellType != nil {
		cell.CellType = *req.CellType
	}
	if req.Source != nil {
		cell.Source = *req.Source
	}
	if req.ExecutionCount != nil {
		cell.ExecutionCount = *req.ExecutionCount
	}

	return m.Repo.UpdateCell(ctx, cell)
}

func (m *CellModule) UpdateCells(ctx context.Context, notebookID uuid.UUID, req *models.UpdateCellsRequest) error {
	m.Logger.Info().
		Str("notebook_id", notebookID.String()).
		Int("delete_count", len(req.CellsToDelete)).
		Int("upsert_count", len(req.CellsToUpsert)).
		Msg("Updating cells in module")
	return m.Repo.UpdateCells(ctx, notebookID, req)
}


func (m *CellModule) DeleteCell(ctx context.Context, id uuid.UUID) error {
	return m.Repo.DeleteCell(ctx, id)
}

func (m *CellModule) CreateCellOutput(ctx context.Context, req *models.CreateCellOutputRequest) (*models.CellOutput, error) {
	if req == nil {
		return nil, errors.New("invalid create cell output request")
	}

	output := &models.CellOutput{
		ID:          uuid.New(),
		CellID:      req.CellID,
		OutputIndex: req.OutputIndex,
		Type:        req.Type,
		DataJSON:    req.DataJSON,
		MinioURL:    req.MinioURL,
	}
	m.Logger.Debug().Str("generated_output_id", output.ID.String()).Str("cell_id", output.CellID.ToUUID().String()).Msg("CellModule: Creating cell output")
	createdOutput, err := m.Repo.CreateCellOutput(ctx, output)
	if err != nil {
		m.Logger.Error().Err(err).Msg("CellModule: Failed to create cell output in repository")
		return nil, err
	}
	m.Logger.Info().Str("created_output_id", createdOutput.ID.String()).Msg("CellModule: Successfully created cell output in repository")
	return createdOutput, nil
}

func (m *CellModule) GetCellOutputsByCellID(ctx context.Context, cellID uuid.UUID) ([]*models.CellOutput, error) {
	return m.Repo.GetCellOutputsByCellID(ctx, cellID)
}

func (m *CellModule) DeleteCellOutput(ctx context.Context, id uuid.UUID) error {
	return m.Repo.DeleteCellOutput(ctx, id)
}