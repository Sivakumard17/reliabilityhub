package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"reliabilityhub.dev/api/internal/config"
	"reliabilityhub.dev/api/internal/repository"
	"reliabilityhub.dev/api/internal/server"
	"reliabilityhub.dev/api/pkg/logger"
)

var Version = "dev"

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString("FATAL: " + err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.MustNew(cfg.Server.Environment)
	defer log.Sync() //nolint:errcheck

	log.Info("starting ReliabilityHub API", zap.String("version", Version))

	ctx := context.Background()
	db, err := repository.NewPool(ctx, cfg.Database, log)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()

	srv := server.New(cfg, log, db)

	serverErr := make(chan error, 1)
	go func() { serverErr <- srv.Start() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatal("server failed", zap.Error(err))
	case sig := <-quit:
		log.Info("shutdown signal", zap.String("signal", sig.String()))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", zap.Error(err))
		os.Exit(1)
	}
}
