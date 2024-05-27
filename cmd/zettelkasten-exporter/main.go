package main

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if len(os.Args) != 2 {
		slog.Error("Wrong arguments")
		os.Exit(1)
	}
	dir := os.Args[1]
	absolute_path, err := filepath.Abs(dir)
	if err != nil {
		slog.Error("Error getting absolute file path", slog.Any("error", err))
		os.Exit(1)
	}
	interval, _ := time.ParseDuration("5m")
	collector := collector.NewCollector(absolute_path, interval)

	slog.Info("Starting metrics collector")
	collector.StartCollecting()
	slog.Info("Metrics collector started")

	addr := "0.0.0.0:6969"
	http.Handle("/metrics", promhttp.Handler())
	slog.Info("Starting HTTP server", slog.String("address", addr))
	err = http.ListenAndServe(addr, nil)
	slog.Info("Error on HTTP server", slog.Any("error", err))
	os.Exit(1)
}
