package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ifaisalabid1/todo-app/internal/config"
	"github.com/ifaisalabid1/todo-app/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	logger := slog.Default()

	logger.Info("starting server", slog.String("env", cfg.Server.Env), slog.String("host", cfg.Server.Host), slog.String("port", cfg.Server.Port))

	pool, err := database.NewPool(cfg.Database, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.ClosePool(pool)
}
