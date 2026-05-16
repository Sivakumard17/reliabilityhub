package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SLO model (kept local to avoid import cycle)
type SLORecord struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ClusterID   *uuid.UUID `json:"cluster_id,omitempty"`
	Service     string     `json:"service"`
	SLOType     string     `json:"slo_type"`
	Target      float64    `json:"target"`
	WindowDays  int        `json:"window_days"`
	PromQLGood  string     `json:"promql_good"`
	PromQLTotal string     `json:"promql_total"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SLOSnapshotRecord struct {
	ID               uuid.UUID
	SLOID            uuid.UUID
	Compliance       float64
	ErrorBudgetTotal float64
	ErrorBudgetUsed  float64
	ErrorBudgetPct   float64
	BurnRate1h       float64
	BurnRate6h       float64
	BurnRate24h      float64
	SnapshotAt       time.Time
}

type CreateSLOInput struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Service     string     `json:"service"`
	SLOType     string     `json:"slo_type"`
	Target      float64    `json:"target"`
	WindowDays  int        `json:"window_days"`
	PromQLGood  string     `json:"promql_good"`
	PromQLTotal string     `json:"promql_total"`
}

type SLORepository struct {
	db *pgxpool.Pool
}

func NewSLORepository(db *pgxpool.Pool) *SLORepository {
	return &SLORepository{db: db}
}

func (r *SLORepository) Create(ctx context.Context, req CreateSLOInput) (*SLORecord, error) {
	query := `
		INSERT INTO slos (
			name, description, service, slo_type,
			target, window_days, promql_good, promql_total, is_active
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,true)
		RETURNING id, name, description, cluster_id, service,
			slo_type, target, window_days,
			promql_good, promql_total, is_active,
			created_at, updated_at`

	row := r.db.QueryRow(ctx, query,
		req.Name, req.Description, req.Service, req.SLOType,
		req.Target, req.WindowDays, req.PromQLGood, req.PromQLTotal,
	)
	return scanSLORecord(row)
}

func (r *SLORepository) GetByID(ctx context.Context, id uuid.UUID) (*SLORecord, error) {
	query := `
		SELECT id, name, description, cluster_id, service,
			slo_type, target, window_days,
			promql_good, promql_total, is_active,
			created_at, updated_at
		FROM slos WHERE id = $1`

	s, err := scanSLORecord(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *SLORepository) ListActive(ctx context.Context) ([]*SLORecord, error) {
	query := `
		SELECT id, name, description, cluster_id, service,
			slo_type, target, window_days,
			promql_good, promql_total, is_active,
			created_at, updated_at
		FROM slos WHERE is_active = true
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*SLORecord{}
	for rows.Next() {
		s, err := scanSLORecord(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *SLORepository) SaveSnapshot(ctx context.Context, snap *SLOSnapshotRecord) error {
	query := `
		INSERT INTO slo_snapshots (
			slo_id, compliance,
			error_budget_total, error_budget_used, error_budget_pct,
			burn_rate_1h, burn_rate_6h, burn_rate_24h,
			snapshot_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())`

	_, err := r.db.Exec(ctx, query,
		snap.SLOID, snap.Compliance,
		snap.ErrorBudgetTotal, snap.ErrorBudgetUsed, snap.ErrorBudgetPct,
		snap.BurnRate1h, snap.BurnRate6h, snap.BurnRate24h,
	)
	return err
}

func (r *SLORepository) GetLatestSnapshot(ctx context.Context, sloID uuid.UUID) (*SLOSnapshotRecord, error) {
	query := `
		SELECT id, slo_id, compliance,
			error_budget_total, error_budget_used, error_budget_pct,
			burn_rate_1h, burn_rate_6h, burn_rate_24h,
			snapshot_at
		FROM slo_snapshots
		WHERE slo_id = $1
		ORDER BY snapshot_at DESC
		LIMIT 1`

	return scanSnapshotRecord(r.db.QueryRow(ctx, query, sloID))
}

type sloRowScanner interface{ Scan(dest ...any) error }

func scanSLORecord(row sloRowScanner) (*SLORecord, error) {
	var s SLORecord
	var clusterID *uuid.UUID
	var promQLGood, promQLTotal *string

	err := row.Scan(
		&s.ID, &s.Name, &s.Description, &clusterID, &s.Service,
		&s.SLOType, &s.Target, &s.WindowDays,
		&promQLGood, &promQLTotal, &s.IsActive,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan slo: %w", err)
	}
	s.ClusterID = clusterID
	if promQLGood != nil {
		s.PromQLGood = *promQLGood
	}
	if promQLTotal != nil {
		s.PromQLTotal = *promQLTotal
	}
	return &s, nil
}

type snapRowScanner interface{ Scan(dest ...any) error }

func scanSnapshotRecord(row snapRowScanner) (*SLOSnapshotRecord, error) {
	var snap SLOSnapshotRecord
	err := row.Scan(
		&snap.ID, &snap.SLOID, &snap.Compliance,
		&snap.ErrorBudgetTotal, &snap.ErrorBudgetUsed, &snap.ErrorBudgetPct,
		&snap.BurnRate1h, &snap.BurnRate6h, &snap.BurnRate24h,
		&snap.SnapshotAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan snapshot: %w", err)
	}
	return &snap, nil
}
