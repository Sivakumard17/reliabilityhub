package incidents

import (
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

type Status string

const (
	StatusOpen          Status = "open"
	StatusInvestigating Status = "investigating"
	StatusIdentified    Status = "identified"
	StatusMonitoring    Status = "monitoring"
	StatusResolved      Status = "resolved"
)

type Incident struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Severity    Severity          `json:"severity"`
	Status      Status            `json:"status"`
	ClusterID   *uuid.UUID        `json:"cluster_id,omitempty"`
	Service     string            `json:"service"`
	AlertName   string            `json:"alert_name"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartedAt   time.Time         `json:"started_at"`
	ResolvedAt  *time.Time        `json:"resolved_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type CreateIncidentRequest struct {
	Title       string            `json:"title"       binding:"required,min=3,max=500"`
	Description string            `json:"description"`
	Severity    Severity          `json:"severity"    binding:"required"`
	ClusterID   *uuid.UUID        `json:"cluster_id"`
	Service     string            `json:"service"`
	AlertName   string            `json:"alert_name"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type UpdateStatusRequest struct {
	Status Status `json:"status" binding:"required"`
}

type ListIncidentFilters struct {
	Status    *Status
	Severity  *Severity
	ClusterID *uuid.UUID
	Service   string
	Limit     int
	Offset    int
}
