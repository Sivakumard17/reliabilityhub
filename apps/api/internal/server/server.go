package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/config"
	"reliabilityhub.dev/api/internal/handler"
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
	s.router.Use(custommiddleware.RequestLogger(s.log))
	s.router.Use(custommiddleware.RequestID())
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (s *Server) setupRoutes() {
	s.router.GET("/healthz", handler.Healthz(s.log))
	s.router.GET("/readyz", handler.Readyz(s.log))
	s.router.GET("/metrics", handler.Metrics())

	v1 := s.router.Group("/api/v1")

	// Incidents
	incidentRepo    := repository.NewIncidentRepository(s.db)
	incidentSvc     := service.NewIncidentService(incidentRepo, s.log)
	incidentHandler := handler.NewIncidentHandler(incidentSvc, s.log)

	inc := v1.Group("/incidents")
	inc.POST("",             incidentHandler.Create)
	inc.GET("",              incidentHandler.List)
	inc.GET("/:id",          incidentHandler.GetByID)
	inc.PATCH("/:id/status", incidentHandler.UpdateStatus)

	// SLOs
	sloRepo    := repository.NewSLORepository(s.db)
	sloHandler := handler.NewSLOHandler(sloRepo, s.log)

	slos := v1.Group("/slos")
	slos.POST("",    sloHandler.Create)
	slos.GET("",     sloHandler.List)
	slos.GET("/:id", sloHandler.GetByID)
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
