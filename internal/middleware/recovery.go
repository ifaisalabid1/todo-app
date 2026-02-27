package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/ifaisalabid1/todo-app/pkg/response"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						slog.String("error", err.(string)),
						slog.String("stack", string(debug.Stack())),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
					)

					response.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "unexpected error occurred")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
