package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"prod-health-check-api/internal/config"
	"prod-health-check-api/internal/db"
	"prod-health-check-api/internal/http/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	mux := http.NewServeMux()

	health := handlers.HealthHandler{DB: conn}

	mux.HandleFunc("/live", health.Live)

	mux.HandleFunc("/ready", health.Ready)

	// Optional: keep /health as alias (useful for simple setups)
	mux.HandleFunc("/health", health.Ready)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(srv, 5*time.Second)
}

func gracefulShutdown(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_ = srv.Shutdown(ctx)
	log.Println("shutdown complete")
}
