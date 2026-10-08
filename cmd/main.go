package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/api"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("could not load .env, falling back to process environment", "error", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn())
	if err != nil {
		slog.Error("could not create database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)

	handler := middleware.Logging(api.NewRouter(queries))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server listening", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// dsn prefers DATABASE_URL, falling back to the individual POSTGRES_* variables
// because .env ships DATABASE_URL with unexpanded shell references.
func dsn() string {
	url := os.Getenv("DATABASE_URL")
	if url != "" && !strings.Contains(url, "${") {
		return url
	}

	host := envOr("POSTGRES_HOST", "localhost")
	port := envOr("POSTGRES_PORT", "5432")
	user := envOr("POSTGRES_USER", "user")
	password := envOr("POSTGRES_PASSWORD", "password")
	name := envOr("POSTGRES_DB", "wms_rasha_shafa")

	return "postgresql://" + user + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable"
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
