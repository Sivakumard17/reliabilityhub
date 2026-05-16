package remediation

import (
	"time"
	"github.com/google/uuid"
)

// TriggerType defines what causes remediation
type TriggerType string

const (
	TriggerAlertName  TriggerType = "alert_name"
	TriggerSeverity   TriggerType = "severity"
	TriggerService    TriggerType = "service"
)

// ActionType defines what remediation does
type ActionType string

const (
	ActionRestartDeployment ActionType = "restart_deployment"
	ActionScaleDeployment   ActionType = "scale_deployment"
	ActionRollbackDeployment ActionType = "rollback_deployment"
	ActionCreateIncident    ActionType = "create_incident"
)

// Policy defines when and how to remediate
type Policy struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	TriggerType TriggerType       `json:"trigger_type"`
	Conditions  PolicyConditions  `json:"conditions"`
	Actions     []PolicyAction    `json:"actions"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
}

// PolicyConditions defines what must match to trigger
type PolicyConditions struct {
	AlertName  string `json:"alert_name,omitempty"`
	Severity   string `json:"severity,omitempty"`
	Service    string `json:"service,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
}

// PolicyAction defines what to do
type PolicyAction struct {
	Type       ActionType        `json:"type"`
	Parameters map[string]string `json:"parameters"`
	// e.g. {"deployment": "payment-service", "namespace": "production"}
	// e.g. {"replicas": "3"}
}

// ExecutionResult records what happened
type ExecutionResult struct {
	PolicyID   uuid.UUID  `json:"policy_id"`
	PolicyName string     `json:"policy_name"`
	Action     ActionType `json:"action"`
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	ExecutedAt time.Time  `json:"executed_at"`
}

// RemediationEvent is written to audit log
type RemediationEvent struct {
	IncidentID uuid.UUID         `json:"incident_id"`
	PolicyID   uuid.UUID         `json:"policy_id"`
	Actions    []ExecutionResult `json:"actions"`
	TriggeredAt time.Time        `json:"triggered_at"`
}
