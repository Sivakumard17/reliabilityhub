package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/incidents"
	"reliabilityhub.dev/api/internal/service"
)

// WebhookHandler handles incoming webhooks from external systems
type WebhookHandler struct {
	incidentSvc *service.IncidentService
	log         *zap.Logger
}

func NewWebhookHandler(svc *service.IncidentService, log *zap.Logger) *WebhookHandler {
	return &WebhookHandler{incidentSvc: svc, log: log}
}

// AlertManager handles POST /api/v1/webhooks/alertmanager
// This is called by AlertManager when alerts fire or resolve
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
		zap.String("receiver", payload.Receiver),
	)

	created := 0
	resolved := 0

	for _, alert := range payload.Alerts {
		switch alert.Status {
		case "firing":
			if err := h.handleFiringAlert(c, alert); err != nil {
				h.log.Error("failed to handle firing alert",
					zap.String("fingerprint", alert.Fingerprint),
					zap.Error(err),
				)
				continue
			}
			created++

		case "resolved":
			h.log.Info("alert resolved",
				zap.String("fingerprint", alert.Fingerprint),
				zap.String("alertname", alert.Labels["alertname"]),
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

func (h *WebhookHandler) handleFiringAlert(c *gin.Context, alert incidents.Alert) error {
	req := incidents.CreateIncidentRequest{
		Title:       incidents.TitleFromAlert(alert),
		Description: alert.Annotations["description"],
		Severity:    incidents.SeverityFromLabels(alert.Labels),
		Service:     incidents.ServiceFromAlert(alert),
		AlertName:   alert.Labels["alertname"],
		Labels:      alert.Labels,
		Annotations: alert.Annotations,
	}

	inc, err := h.incidentSvc.Create(c.Request.Context(), req)
	if err != nil {
		return err
	}

	h.log.Info("incident created from alert",
		zap.String("incident_id", inc.ID.String()),
		zap.String("title", inc.Title),
		zap.String("severity", string(inc.Severity)),
		zap.String("fingerprint", alert.Fingerprint),
	)

	return nil
}

// Generic webhook for future integrations
// POST /api/v1/webhooks/generic
func (h *WebhookHandler) Generic(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid payload"})
		return
	}

	h.log.Info("generic webhook received",
		zap.Any("payload", body),
	)

	c.JSON(http.StatusOK, APIResponse{Data: gin.H{"status": "received"}})
}
