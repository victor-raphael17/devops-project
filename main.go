package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	projectName = "Projeto Korp"
	defaultPort = "8080"
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
	router := http.NewServeMux()

	router.Handle("GET /projeto-korp", promhttp.InstrumentHandlerCounter(
		requestsTotal, http.HandlerFunc(handleProjetoKorp)))
	
	router.Handle("GET /metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              ":" + port(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
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

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}

	return defaultPort
}
