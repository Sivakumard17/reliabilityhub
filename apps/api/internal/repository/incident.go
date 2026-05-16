package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"reliabilityhub.dev/api/internal/incidents"
)

var ErrNotFound = errors.New("record not found")

type IncidentRepository struct {
	db *pgxpool.Pool
}

func NewIncidentRepository(db *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) Create(ctx context.Context, req incidents.CreateIncidentRequest) (*incidents.Incident, error) {
	labelsJSON, err := json.Marshal(req.Labels)
	if err != nil {
		return nil, fmt.Errorf("marshal labels: %w", err)
	}
	annotationsJSON, err := json.Marshal(req.Annotations)
	if err != nil {
		return nil, fmt.Errorf("marshal annotations: %w", err)
	}

	query := `
		INSERT INTO incidents (
			title, description, severity, status,
			cluster_id, service, alert_name,
			labels, annotations, started_at
		) VALUES ($1, $2, $3, 'open', $4, $5, $6, $7, $8, NOW())
		RETURNING
			id, title, description, severity, status,
			cluster_id, service, alert_name,
			labels, annotations,
			started_at, resolved_at, created_at, updated_at`

	row := r.db.QueryRow(ctx, query,
		req.Title, req.Description, req.Severity,
		req.ClusterID, req.Service, req.AlertName,
		labelsJSON, annotationsJSON,
	)
	return scanIncident(row)
}

func (r *IncidentRepository) GetByID(ctx context.Context, id uuid.UUID) (*incidents.Incident, error) {
	query := `
		SELECT id, title, description, severity, status,
			cluster_id, service, alert_name,
			labels, annotations,
			started_at, resolved_at, created_at, updated_at
		FROM incidents WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	incident, err := scanIncident(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return incident, nil
}

func (r *IncidentRepository) List(ctx context.Context, filters incidents.ListIncidentFilters) ([]*incidents.Incident, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if filters.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, *filters.Status)
		i++
	}
	if filters.Severity != nil {
		where += fmt.Sprintf(" AND severity = $%d", i)
		args = append(args, *filters.Severity)
		i++
	}
	if filters.Service != "" {
		where += fmt.Sprintf(" AND service = $%d", i)
		args = append(args, filters.Service)
		i++
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM incidents %s", where)
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	limit := filters.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	listQuery := fmt.Sprintf(`
		SELECT id, title, description, severity, status,
			cluster_id, service, alert_name,
			labels, annotations,
			started_at, resolved_at, created_at, updated_at
		FROM incidents %s
		ORDER BY started_at DESC
		LIMIT $%d OFFSET $%d`, where, i, i+1)

	rows, err := r.db.Query(ctx, listQuery, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list: %w", err)
	}
	defer rows.Close()

	result := make([]*incidents.Incident, 0)
	for rows.Next() {
		inc, err := scanIncident(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, inc)
	}
	return result, total, nil
}

func (r *IncidentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status incidents.Status) (*incidents.Incident, error) {
	var resolvedAt *time.Time
	if status == incidents.StatusResolved {
		now := time.Now()
		resolvedAt = &now
	}

	query := `
		UPDATE incidents
		SET status = $1, resolved_at = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING
			id, title, description, severity, status,
			cluster_id, service, alert_name,
			labels, annotations,
			started_at, resolved_at, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, status, resolvedAt, id)
	inc, err := scanIncident(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return inc, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanIncident(row scanner) (*incidents.Incident, error) {
	var inc incidents.Incident
	var labelsJSON, annotationsJSON []byte
	var clusterID *uuid.UUID

	err := row.Scan(
		&inc.ID, &inc.Title, &inc.Description,
		&inc.Severity, &inc.Status,
		&clusterID, &inc.Service, &inc.AlertName,
		&labelsJSON, &annotationsJSON,
		&inc.StartedAt, &inc.ResolvedAt,
		&inc.CreatedAt, &inc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	inc.ClusterID = clusterID

	if err := json.Unmarshal(labelsJSON, &inc.Labels); err != nil {
		inc.Labels = map[string]string{}
	}
	if err := json.Unmarshal(annotationsJSON, &inc.Annotations); err != nil {
		inc.Annotations = map[string]string{}
	}
	return &inc, nil
}
