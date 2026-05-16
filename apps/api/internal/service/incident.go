package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/incidents"
	"reliabilityhub.dev/api/internal/repository"
)

type IncidentService struct {
	repo *repository.IncidentRepository
	log  *zap.Logger
}

func NewIncidentService(repo *repository.IncidentRepository, log *zap.Logger) *IncidentService {
	return &IncidentService{repo: repo, log: log}
}

func (s *IncidentService) Create(ctx context.Context, req incidents.CreateIncidentRequest) (*incidents.Incident, error) {
	if req.Labels == nil {
		req.Labels = map[string]string{}
	}
	if req.Annotations == nil {
		req.Annotations = map[string]string{}
	}

	inc, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create incident: %w", err)
	}

	s.log.Info("incident created",
		zap.String("id", inc.ID.String()),
		zap.String("severity", string(inc.Severity)),
		zap.String("title", inc.Title),
	)
	return inc, nil
}

func (s *IncidentService) GetByID(ctx context.Context, id uuid.UUID) (*incidents.Incident, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *IncidentService) List(ctx context.Context, filters incidents.ListIncidentFilters) ([]*incidents.Incident, int, error) {
	return s.repo.List(ctx, filters)
}

func (s *IncidentService) UpdateStatus(ctx context.Context, id uuid.UUID, status incidents.Status) (*incidents.Incident, error) {
	inc, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}

	s.log.Info("incident status updated",
		zap.String("id", id.String()),
		zap.String("status", string(status)),
	)
	return inc, nil
}
