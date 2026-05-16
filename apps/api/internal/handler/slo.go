package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/repository"
)

type SLOHandler struct {
	repo *repository.SLORepository
	log  *zap.Logger
}

func NewSLOHandler(repo *repository.SLORepository, log *zap.Logger) *SLOHandler {
	return &SLOHandler{repo: repo, log: log}
}

type CreateSLORequest struct {
	Name        string  `json:"name"         binding:"required"`
	Description string  `json:"description"`
	Service     string  `json:"service"      binding:"required"`
	SLOType     string  `json:"slo_type"     binding:"required"`
	Target      float64 `json:"target"       binding:"required"`
	WindowDays  int     `json:"window_days"  binding:"required"`
	PromQLGood  string  `json:"promql_good"`
	PromQLTotal string  `json:"promql_total"`
}

type SLOStatusResponse struct {
	SLO      *repository.SLORecord         `json:"slo"`
	Snapshot *repository.SLOSnapshotRecord `json:"snapshot,omitempty"`
	Status   string                        `json:"status"`
}

func burnRateStatus(snap *repository.SLOSnapshotRecord) string {
	if snap == nil {
		return "unknown"
	}
	switch {
	case snap.BurnRate1h > 14.4:
		return "critical"
	case snap.BurnRate6h > 6:
		return "warning"
	case snap.BurnRate1h > 1:
		return "degraded"
	default:
		return "healthy"
	}
}

func (h *SLOHandler) Create(c *gin.Context) {
	var req CreateSLORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: err.Error()})
		return
	}

	s, err := h.repo.Create(c.Request.Context(), repository.CreateSLOInput{
		Name:        req.Name,
		Description: req.Description,
		Service:     req.Service,
		SLOType:     req.SLOType,
		Target:      req.Target,
		WindowDays:  req.WindowDays,
		PromQLGood:  req.PromQLGood,
		PromQLTotal: req.PromQLTotal,
	})
	if err != nil {
		h.log.Error("create slo failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to create SLO"})
		return
	}
	c.JSON(http.StatusCreated, APIResponse{Data: s})
}

func (h *SLOHandler) List(c *gin.Context) {
	slos, err := h.repo.ListActive(c.Request.Context())
	if err != nil {
		h.log.Error("list slos failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to list SLOs"})
		return
	}

	statuses := make([]SLOStatusResponse, 0, len(slos))
	for _, s := range slos {
		snap, _ := h.repo.GetLatestSnapshot(c.Request.Context(), s.ID)
		statuses = append(statuses, SLOStatusResponse{
			SLO:      s,
			Snapshot: snap,
			Status:   burnRateStatus(snap),
		})
	}

	total := len(statuses)
	c.JSON(http.StatusOK, APIResponse{Data: statuses, Total: &total})
}

func (h *SLOHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid SLO ID"})
		return
	}

	s, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{Error: "SLO not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to get SLO"})
		return
	}

	snap, _ := h.repo.GetLatestSnapshot(c.Request.Context(), id)
	c.JSON(http.StatusOK, APIResponse{
		Data: SLOStatusResponse{SLO: s, Snapshot: snap, Status: burnRateStatus(snap)},
	})
}
