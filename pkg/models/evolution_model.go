package models

import (
	"time"

	"github.com/google/uuid"
)

// EvolutionRun represents an evolution run for a cell.
type EvolutionRun struct {
	ID           uuid.UUID       `json:"id"`
	SourceCellID StringUUID      `json:"source_cell_id"`
	StartTime    time.Time       `json:"start_time"`
	EndTime      *time.Time      `json:"end_time,omitempty"`
	Status       string          `json:"status"`
	Variations   []CellVariation `json:"variations,omitempty"`
}

// CellVariation represents a single variation within an evolution run.
type CellVariation struct {
	ID              uuid.UUID  `json:"id"`
	EvolutionRunID  uuid.UUID  `json:"evolution_run_id"`
	Code            string     `json:"code"`
	Metric          float64    `json:"metric"`
	IsBest          bool       `json:"is_best"`
	Generation      int        `json:"generation"`
	ParentVariantID *uuid.UUID `json:"parent_variant_id,omitempty"`
}