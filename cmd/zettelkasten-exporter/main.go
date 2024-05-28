package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Loaded config", slog.Any("config", cfg))

	absolute_path, err := filepath.Abs(cfg.ZettelkastenDirectory)
	collector := collector.NewCollector(absolute_path, cfg.ScrapeInterval)

	slog.Info("Starting metrics collector")
	collector.StartCollecting()
	slog.Info("Metrics collector started")

	addr := fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	http.Handle("/metrics", promhttp.Handler())
	slog.Info("Starting HTTP server", slog.String("address", addr))
	err = http.ListenAndServe(addr, nil)
	slog.Info("Error on HTTP server", slog.Any("error", err))
	os.Exit(1)
}
