package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/incidents"
	"reliabilityhub.dev/api/internal/remediation"
	"reliabilityhub.dev/api/internal/repository"
)

type RemediationHandler struct {
	engine       *remediation.Engine
	incidentRepo *repository.IncidentRepository
	log          *zap.Logger
}

func NewRemediationHandler(
	engine *remediation.Engine,
	repo *repository.IncidentRepository,
	log *zap.Logger,
) *RemediationHandler {
	return &RemediationHandler{
		engine:       engine,
		incidentRepo: repo,
		log:          log,
	}
}

// Trigger manually triggers remediation for an incident
// POST /api/v1/remediation/trigger/:incident_id
func (h *RemediationHandler) Trigger(c *gin.Context) {
	id, err := uuid.Parse(c.Param("incident_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid incident ID"})
		return
	}

	incident, err := h.incidentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{Error: "incident not found"})
		return
	}

	h.log.Info("manual remediation triggered",
		zap.String("incident_id", id.String()),
		zap.String("title", incident.Title),
	)

	results := h.engine.Evaluate(c.Request.Context(), incident)

	if len(results) == 0 {
		c.JSON(http.StatusOK, APIResponse{
			Data: gin.H{
				"message":     "no matching remediation policies found",
				"incident_id": id,
			},
		})
		return
	}

	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{
			"incident_id":       id,
			"policies_matched":  len(results),
			"actions_succeeded": successCount,
			"actions_failed":    len(results) - successCount,
			"results":           results,
		},
	})
}

// Policies lists all remediation policies
// GET /api/v1/remediation/policies
func (h *RemediationHandler) Policies(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{
			"policies": []gin.H{
				{
					"name":        "Auto-restart on CrashLoopBackOff",
					"trigger":     "alert_name=KubePodCrashLooping",
					"action":      "restart_deployment",
					"is_active":   true,
				},
				{
					"name":      "Scale up on high load",
					"trigger":   "alert_name=HighCPUUsage",
					"action":    "scale_deployment(replicas=4)",
					"is_active": true,
				},
				{
					"name":      "Restart on high error rate",
					"trigger":   "alert_name=HighErrorRate,severity=critical",
					"action":    "restart_deployment(reliabilityhub-api)",
					"is_active": true,
				},
			},
		},
	})
}

// TriggerAutoRemediation runs remediation asynchronously after incident creation
func TriggerAutoRemediation(
	engine *remediation.Engine,
	incident *incidents.Incident,
	log *zap.Logger,
) {
	go func() {
		ctx := context.Background()
		results := engine.Evaluate(ctx, incident)

		if len(results) > 0 {
			successCount := 0
			for _, r := range results {
				if r.Success {
					successCount++
				}
			}
			log.Info("auto-remediation completed",
				zap.String("incident_id", incident.ID.String()),
				zap.Int("actions_taken", len(results)),
				zap.Int("succeeded", successCount),
			)
		}
	}()
}
