package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/incidents"
	"reliabilityhub.dev/api/internal/remediation"
	"reliabilityhub.dev/api/internal/service"
)

type WebhookHandler struct {
	incidentSvc      *service.IncidentService
	remediationEngine *remediation.Engine
	log              *zap.Logger
}

func NewWebhookHandler(
	svc *service.IncidentService,
	log *zap.Logger,
	engine *remediation.Engine,
) *WebhookHandler {
	return &WebhookHandler{
		incidentSvc:      svc,
		remediationEngine: engine,
		log:              log,
	}
}

func (h *WebhookHandler) AlertManager(c *gin.Context) {
	var payload incidents.AlertManagerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.log.Error("invalid alertmanager payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid payload"})
		return
	}

	h.log.Info("alertmanager webhook received",
		zap.String("status", payload.Status),
		zap.Int("alert_count", len(payload.Alerts)),
	)

	created := 0
	resolved := 0

	for _, alert := range payload.Alerts {
		switch alert.Status {
		case "firing":
			inc, err := h.incidentSvc.Create(c.Request.Context(),
				incidents.CreateIncidentRequest{
					Title:       incidents.TitleFromAlert(alert),
					Description: alert.Annotations["description"],
					Severity:    incidents.SeverityFromLabels(alert.Labels),
					Service:     incidents.ServiceFromAlert(alert),
					AlertName:   alert.Labels["alertname"],
					Labels:      alert.Labels,
					Annotations: alert.Annotations,
				},
			)
			if err != nil {
				h.log.Error("failed to create incident from alert", zap.Error(err))
				continue
			}
			created++

			// Auto-remediate if engine is available
			if h.remediationEngine != nil {
				TriggerAutoRemediation(h.remediationEngine, inc, h.log)
			}

		case "resolved":
			h.log.Info("alert resolved",
				zap.String("fingerprint", alert.Fingerprint),
			)
			resolved++
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{
			"received": len(payload.Alerts),
			"created":  created,
			"resolved": resolved,
		},
	})
}

func (h *WebhookHandler) Generic(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid payload"})
		return
	}
	h.log.Info("generic webhook received", zap.Any("payload", body))
	c.JSON(http.StatusOK, APIResponse{Data: gin.H{"status": "received"}})
}
