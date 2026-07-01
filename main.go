package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	projectName     = "Projeto Korp"
	defaultPort     = "8080"
	shutdownTimeout = 5 * time.Second
)

type response struct {
	Name string `json:"nome"`
	Time string `json:"horario"`
}

var requestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests served by /projeto-korp.",
	},
	[]string{"code", "method"},
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:              ":" + port(),
		Handler:           newRouter(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() { serverErr <- server.ListenAndServe() }()

	log.Printf("server listening on %s", server.Addr)

	select {
	case err := <-serverErr:
		log.Fatalf("server stopped: %v", err)
	case <-ctx.Done():
		log.Print("shutdown signal received, draining connections")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Print("server stopped")
}

func newRouter() http.Handler {
	router := http.NewServeMux()

	router.Handle("GET /projeto-korp", promhttp.InstrumentHandlerCounter(
		requestsTotal, http.HandlerFunc(handleProjetoKorp)))

	// Not instrumented, so container healthchecks don't inflate the
	// request-volume metrics.
	router.HandleFunc("GET /healthz", handleHealthz)

	router.Handle("GET /metrics", promhttp.Handler())

	return router
}

func handleProjetoKorp(w http.ResponseWriter, r *http.Request) {
	body := response{
		Name: projectName,
		Time: time.Now().UTC().Format(time.RFC3339),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "failed to build response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(payload)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}

	return defaultPort
}
