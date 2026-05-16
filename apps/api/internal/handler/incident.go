package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/incidents"
	"reliabilityhub.dev/api/internal/repository"
	"reliabilityhub.dev/api/internal/service"
)

type IncidentHandler struct {
	svc *service.IncidentService
	log *zap.Logger
}

func NewIncidentHandler(svc *service.IncidentService, log *zap.Logger) *IncidentHandler {
	return &IncidentHandler{svc: svc, log: log}
}

type APIResponse struct {
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Total   *int   `json:"total,omitempty"`
	Page    *int   `json:"page,omitempty"`
	PerPage *int   `json:"per_page,omitempty"`
}

func (h *IncidentHandler) Create(c *gin.Context) {
	var req incidents.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: err.Error()})
		return
	}

	inc, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		h.log.Error("create incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to create incident"})
		return
	}
	c.JSON(http.StatusCreated, APIResponse{Data: inc})
}

func (h *IncidentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid incident ID"})
		return
	}

	inc, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{Error: "incident not found"})
			return
		}
		h.log.Error("get incident failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to get incident"})
		return
	}
	c.JSON(http.StatusOK, APIResponse{Data: inc})
}

func (h *IncidentHandler) List(c *gin.Context) {
	filters := incidents.ListIncidentFilters{}

	if s := c.Query("status"); s != "" {
		st := incidents.Status(s)
		filters.Status = &st
	}
	if s := c.Query("severity"); s != "" {
		sv := incidents.Severity(s)
		filters.Severity = &sv
	}
	if s := c.Query("service"); s != "" {
		filters.Service = s
	}

	page, perPage := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if p := c.Query("per_page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 && v <= 100 {
			perPage = v
		}
	}
	filters.Limit = perPage
	filters.Offset = (page - 1) * perPage

	result, total, err := h.svc.List(c.Request.Context(), filters)
	if err != nil {
		h.log.Error("list incidents failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to list incidents"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data:    result,
		Total:   &total,
		Page:    &page,
		PerPage: &perPage,
	})
}

func (h *IncidentHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "invalid incident ID"})
		return
	}

	var body incidents.UpdateStatusRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: err.Error()})
		return
	}

	inc, err := h.svc.UpdateStatus(c.Request.Context(), id, body.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{Error: "incident not found"})
			return
		}
		h.log.Error("update status failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "failed to update incident"})
		return
	}
	c.JSON(http.StatusOK, APIResponse{Data: inc})
}
