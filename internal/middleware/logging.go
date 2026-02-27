package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New().String()

			ctx := r.Context()
			ctx = context.WithValue(ctx, "request_id", requestID)
			r = r.WithContext(ctx)

			w.Header().Set("X-REQUEST-ID", requestID)

			next.ServeHTTP(w, r)

			duration := time.Since(start)

			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", requestID),
				slog.String("ip", r.RemoteAddr),
				slog.Duration("duration", duration),
			)
		})
	}
}
