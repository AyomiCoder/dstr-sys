package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/AyomiCoder/dstr-sys.git/internal/config"
	httpserver "github.com/AyomiCoder/dstr-sys.git/internal/http"
	"github.com/AyomiCoder/dstr-sys.git/internal/notification"
)

func main() {
	cfg := config.Load()
	store := notification.NewStore()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpserver.NewRouter(cfg, store),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("starting %s on %s", cfg.ServiceName, server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
