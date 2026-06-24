// Package main runs a small HTTP server with a single endpoint,
// GET /projeto-korp, that returns the project name and the current UTC time.
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

// response is the JSON body returned by GET /projeto-korp.
type response struct {
	Name string `json:"nome"`
	Time string `json:"horario"`
}

// requestsTotal counts the requests served by GET /projeto-korp, labelled by
// HTTP status code and method, so Prometheus can measure request volume.
var requestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests served by /projeto-korp.",
	},
	[]string{"code", "method"},
)

func main() {
	// Register the endpoints on our own router (instead of the global one).
	router := http.NewServeMux()
	// InstrumentHandlerCounter increments requestsTotal on every request,
	// keeping the metric logic out of the business handler.
	router.Handle("GET /projeto-korp", promhttp.InstrumentHandlerCounter(
		requestsTotal, http.HandlerFunc(handleProjetoKorp)))
	// /metrics exposes the collected metrics in the Prometheus exposition format.
	router.Handle("GET /metrics", promhttp.Handler())

	// ReadHeaderTimeout guards against slow clients holding the connection open.
	server := &http.Server{
		Addr:              ":" + port(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ListenAndServe blocks until the server is stopped or hits an error.
	log.Printf("server listening on %s", server.Addr)

	if listenError := server.ListenAndServe(); listenError != nil {
		log.Fatalf("server stopped: %v", listenError)
	}
}

// handleProjetoKorp returns the project name and the current UTC time,
// recomputed on every request.
func handleProjetoKorp(responseWriter http.ResponseWriter, request *http.Request) {
	body := response{
		Name: projectName,
		Time: time.Now().UTC().Format(time.RFC3339),
	}

	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	if encodeError := json.NewEncoder(responseWriter).Encode(body); encodeError != nil {
		http.Error(responseWriter, "failed to build response", http.StatusInternalServerError)
	}
}

// port returns the value of the PORT environment variable,
// or the default port when PORT is not set.
func port() string {
	if portFromEnvironment := os.Getenv("PORT"); portFromEnvironment != "" {
		return portFromEnvironment
	}

	return defaultPort
}
