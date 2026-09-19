package middleware

import (
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/ctxlog"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		reqLogger := slog.Default().With(
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		ctx := ctxlog.WithLogger(r.Context(), reqLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
