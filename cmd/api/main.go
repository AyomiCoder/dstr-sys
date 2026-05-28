package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/AyomiCoder/dstr-sys.git/internal/config"
	"github.com/AyomiCoder/dstr-sys.git/internal/db/postgres"
	httpserver "github.com/AyomiCoder/dstr-sys.git/internal/http"
	"github.com/AyomiCoder/dstr-sys.git/internal/notification"
)

func main() {
	if err := config.LoadDotEnvIfPresent(".env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to initialize postgres pool: %v", err)
	}
	defer pool.Close()

	repo := notification.NewPostgresRepository(pool)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpserver.NewRouter(cfg, repo),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("starting %s on %s (postgres connected)", cfg.ServiceName, server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
