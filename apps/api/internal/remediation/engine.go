package remediation

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/incidents"
	"reliabilityhub.dev/api/internal/k8sclient"
)

// Engine evaluates remediation policies against incidents
// and executes the appropriate actions
type Engine struct {
	k8s      *k8sclient.Client
	log      *zap.Logger
	policies []Policy
}

func NewEngine(k8s *k8sclient.Client, log *zap.Logger) *Engine {
	e := &Engine{k8s: k8s, log: log}
	e.loadDefaultPolicies()
	return e
}

// loadDefaultPolicies loads built-in remediation policies.
// In production these would be loaded from PostgreSQL.
func (e *Engine) loadDefaultPolicies() {
	e.policies = []Policy{
		{
			Name:        "Auto-restart on CrashLoopBackOff",
			Description: "Restart deployment when pod is crash looping",
			TriggerType: TriggerAlertName,
			Conditions: PolicyConditions{
				AlertName: "KubePodCrashLooping",
			},
			Actions: []PolicyAction{
				{
					Type: ActionRestartDeployment,
					Parameters: map[string]string{
						"namespace": "reliabilityhub-system",
					},
				},
			},
			IsActive: true,
		},
		{
			Name:        "Scale up on high load",
			Description: "Scale deployment when CPU > 80%",
			TriggerType: TriggerAlertName,
			Conditions: PolicyConditions{
				AlertName: "HighCPUUsage",
			},
			Actions: []PolicyAction{
				{
					Type: ActionScaleDeployment,
					Parameters: map[string]string{
						"namespace": "reliabilityhub-system",
						"replicas":  "4",
					},
				},
			},
			IsActive: true,
		},
		{
			Name:        "Restart on high error rate",
			Description: "Restart API on sustained high error rate",
			TriggerType: TriggerAlertName,
			Conditions: PolicyConditions{
				AlertName: "HighErrorRate",
				Severity:  "critical",
			},
			Actions: []PolicyAction{
				{
					Type: ActionRestartDeployment,
					Parameters: map[string]string{
						"namespace":  "reliabilityhub-system",
						"deployment": "reliabilityhub-api",
					},
				},
			},
			IsActive: true,
		},
	}

	e.log.Info("remediation policies loaded",
		zap.Int("count", len(e.policies)),
	)
}

// Evaluate checks if any policies match the incident
// and executes their actions
func (e *Engine) Evaluate(ctx context.Context, incident *incidents.Incident) []ExecutionResult {
	results := []ExecutionResult{}

	e.log.Info("evaluating remediation policies",
		zap.String("incident_id", incident.ID.String()),
		zap.String("alert_name", incident.AlertName),
		zap.String("severity", string(incident.Severity)),
	)

	for _, policy := range e.policies {
		if !policy.IsActive {
			continue
		}

		if !e.matches(policy, incident) {
			continue
		}

		e.log.Info("policy matched",
			zap.String("policy", policy.Name),
			zap.String("incident", incident.Title),
		)

		for _, action := range policy.Actions {
			result := e.execute(ctx, policy, action, incident)
			results = append(results, result)
		}
	}

	return results
}

// matches checks if a policy applies to an incident
func (e *Engine) matches(policy Policy, incident *incidents.Incident) bool {
	c := policy.Conditions

	if c.AlertName != "" && c.AlertName != incident.AlertName {
		return false
	}
	if c.Severity != "" && c.Severity != string(incident.Severity) {
		return false
	}
	if c.Service != "" && c.Service != incident.Service {
		return false
	}
	return true
}

// execute runs a single remediation action
func (e *Engine) execute(
	ctx context.Context,
	policy Policy,
	action PolicyAction,
	incident *incidents.Incident,
) ExecutionResult {
	result := ExecutionResult{
		PolicyName: policy.Name,
		Action:     action.Type,
		ExecutedAt: time.Now(),
	}

	e.log.Info("executing remediation action",
		zap.String("action", string(action.Type)),
		zap.String("policy", policy.Name),
		zap.Any("params", action.Parameters),
	)

	var err error

	switch action.Type {
	case ActionRestartDeployment:
		err = e.executeRestart(ctx, action, incident)

	case ActionScaleDeployment:
		err = e.executeScale(ctx, action, incident)

	case ActionRollbackDeployment:
		err = e.executeRollback(ctx, action, incident)

	default:
		err = fmt.Errorf("unknown action type: %s", action.Type)
	}

	if err != nil {
		result.Success = false
		result.Message = err.Error()
		e.log.Error("remediation action failed",
			zap.String("action", string(action.Type)),
			zap.Error(err),
		)
	} else {
		result.Success = true
		result.Message = fmt.Sprintf("action %s completed successfully", action.Type)
		e.log.Info("remediation action succeeded",
			zap.String("action", string(action.Type)),
			zap.String("policy", policy.Name),
		)
	}

	return result
}

func (e *Engine) executeRestart(
	ctx context.Context,
	action PolicyAction,
	incident *incidents.Incident,
) error {
	namespace := action.Parameters["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	// Use deployment from params, fall back to service name
	deployment := action.Parameters["deployment"]
	if deployment == "" {
		deployment = incident.Service
	}

	if deployment == "" {
		return fmt.Errorf("no deployment specified for restart action")
	}

	return e.k8s.RestartDeployment(ctx, namespace, deployment)
}

func (e *Engine) executeScale(
	ctx context.Context,
	action PolicyAction,
	incident *incidents.Incident,
) error {
	namespace := action.Parameters["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	deployment := action.Parameters["deployment"]
	if deployment == "" {
		deployment = incident.Service
	}

	replicasStr := action.Parameters["replicas"]
	if replicasStr == "" {
		return fmt.Errorf("replicas parameter required for scale action")
	}

	replicas, err := strconv.ParseInt(replicasStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid replicas value: %s", replicasStr)
	}

	return e.k8s.ScaleDeployment(ctx, namespace, deployment, int32(replicas))
}

func (e *Engine) executeRollback(
	ctx context.Context,
	action PolicyAction,
	incident *incidents.Incident,
) error {
	namespace := action.Parameters["namespace"]
	if namespace == "" {
		namespace = "default"
	}

	deployment := action.Parameters["deployment"]
	if deployment == "" {
		deployment = incident.Service
	}

	return e.k8s.RollbackDeployment(ctx, namespace, deployment)
}
