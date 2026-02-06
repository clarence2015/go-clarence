package main

import (
	"log"
	"os"

	"github.com/clarence/GoBlog/internal/config"
	"github.com/clarence/GoBlog/internal/infra/db"
	"github.com/clarence/GoBlog/internal/infra/logging"
	"github.com/clarence/GoBlog/internal/transport/httpserver"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logging.NewLogger(cfg.Env)
	logger.Info("starting GoBlog service",
		"env", cfg.Env,
		"http_port", cfg.HTTPPort,
	)

	gormDB, err := db.NewPostgres(cfg.DatabaseDSN, logger)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	engine := httpserver.NewRouter(cfg.Env, logger, gormDB)

	if err := engine.Run(":" + cfg.HTTPPort); err != nil {
		logger.Error("http server exited with error", "err", err)
		os.Exit(1)
	}
}

