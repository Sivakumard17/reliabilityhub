package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/config"
	"reliabilityhub.dev/api/internal/handler"
	"reliabilityhub.dev/api/internal/k8sclient"
	"reliabilityhub.dev/api/internal/remediation"
	"reliabilityhub.dev/api/internal/repository"
	"reliabilityhub.dev/api/internal/service"
	custommiddleware "reliabilityhub.dev/api/pkg/middleware"
)

type Server struct {
	cfg    *config.Config
	log    *zap.Logger
	db     *pgxpool.Pool
	router *gin.Engine
	http   *http.Server
}

func New(cfg *config.Config, log *zap.Logger, db *pgxpool.Pool) *Server {
	if cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	s := &Server{cfg: cfg, log: log, db: db, router: router}
	s.setupMiddleware()
	s.setupRoutes()
	s.http = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(gin.RecoveryWithWriter(nil, func(c *gin.Context, err any) {
		s.log.Error("panic recovered", zap.Any("error", err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}))
	s.router.Use(custommiddleware.AllowAllCORS())
	s.router.Use(custommiddleware.RequestLogger(s.log))
	s.router.Use(custommiddleware.RequestID())
}

func (s *Server) setupRoutes() {
	s.router.GET("/healthz", handler.Healthz(s.log))
	s.router.GET("/readyz", handler.Readyz(s.log))
	s.router.GET("/metrics", handler.Metrics())

	// ── Repositories ──────────────────────────────────────────────────
	incidentRepo := repository.NewIncidentRepository(s.db)
	sloRepo      := repository.NewSLORepository(s.db)

	// ── Services ──────────────────────────────────────────────────────
	incidentSvc := service.NewIncidentService(incidentRepo, s.log)

	// ── Kubernetes client (optional — graceful degradation) ───────────
	kubeconfigPath := os.Getenv("KUBECONFIG")
	k8sClient, err := k8sclient.New(kubeconfigPath, s.log)
	if err != nil {
		s.log.Warn("kubernetes client unavailable — remediation disabled",
			zap.Error(err),
		)
		k8sClient = nil
	}

	// ── Remediation engine ────────────────────────────────────────────
	var remediationEngine *remediation.Engine
	if k8sClient != nil {
		remediationEngine = remediation.NewEngine(k8sClient, s.log)
	}

	// ── Handlers ──────────────────────────────────────────────────────
	incidentHandler    := handler.NewIncidentHandler(incidentSvc, s.log)
	webhookHandler     := handler.NewWebhookHandler(incidentSvc, s.log, remediationEngine)
	sloHandler         := handler.NewSLOHandler(sloRepo, s.log)

	v1 := s.router.Group("/api/v1")

	// Incidents
	inc := v1.Group("/incidents")
	inc.POST("",             incidentHandler.Create)
	inc.GET("",              incidentHandler.List)
	inc.GET("/:id",          incidentHandler.GetByID)
	inc.PATCH("/:id/status", incidentHandler.UpdateStatus)

	// SLOs
	slos := v1.Group("/slos")
	slos.POST("",    sloHandler.Create)
	slos.GET("",     sloHandler.List)
	slos.GET("/:id", sloHandler.GetByID)

	// Webhooks
	webhooks := v1.Group("/webhooks")
	webhooks.POST("/alertmanager", webhookHandler.AlertManager)
	webhooks.POST("/generic",      webhookHandler.Generic)

	// Remediation
	if remediationEngine != nil {
		remHandler := handler.NewRemediationHandler(
			remediationEngine, incidentRepo, s.log,
		)
		rem := v1.Group("/remediation")
		rem.GET("/policies",          remHandler.Policies)
		rem.POST("/trigger/:incident_id", remHandler.Trigger)
	}
}

func (s *Server) Start() error {
	s.log.Info("server starting",
		zap.String("addr", s.http.Addr),
		zap.String("env", s.cfg.Server.Environment),
	)
	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("shutting down...")
	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	s.log.Info("stopped")
	return nil
}
