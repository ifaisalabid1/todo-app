package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ifaisalabid1/todo-app/internal/config"
	"github.com/ifaisalabid1/todo-app/internal/database"
	"github.com/ifaisalabid1/todo-app/internal/handler"
	"github.com/ifaisalabid1/todo-app/internal/repository"
	"github.com/ifaisalabid1/todo-app/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg)

	logger.Info("starting server", slog.String("env", cfg.Server.Env), slog.String("host", cfg.Server.Host), slog.String("port", cfg.Server.Port))

	pool, err := database.NewPool(cfg.Database, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer database.ClosePool(pool)

	todoRepo := repository.NewTodoRepository(pool)
	todoService := service.NewTodoService(todoRepo, logger)
	todoHandler := handler.NewTodoHandler(todoService)

	router := setupRoutes(todoHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("server listening on ", slog.String("addr", srv.Addr))
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	case sig := <-shutdown:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", slog.String("error", err.Error()))
			srv.Close()
			os.Exit(1)
		}

		logger.Info("server stopped gracefully")

	}
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	if cfg.Log.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: getSlogLevel(cfg.Log.Level),
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: getSlogLevel(cfg.Log.Level),
		})
	}

	return slog.New(handler)
}

func getSlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func setupRoutes(todoHandler *handler.TodoHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(time.Minute))

	now := time.Now().UTC().Format(time.RFC3339)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `{"status": "healthy", "timestamp": %s}`, now)
	})

	r.Route("/api/v1", func(r chi.Router) {
		todoHandler.RegisterRoutes(r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"success": "false", "message": "endpoint not found", "timestamp": %s}`, now)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"success": "false", "message": "method not allowed", "timestamp": %s}`, now)
	})

	return r
}
