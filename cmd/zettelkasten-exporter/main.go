package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
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
	if err != nil {
		slog.Error("Error getting absolute path", slog.Any("error", err), slog.String("path", cfg.ZettelkastenDirectory))
		os.Exit(1)
	}
	_, err = os.Stat(absolute_path)
	if err != nil {
		slog.Error("Cannot stat zettelkasten directory", slog.Any("error", err), slog.String("path", absolute_path))
		os.Exit(1)
	}

	fs := os.DirFS(absolute_path)
	collector := collector.NewCollector(fs)

	promHandler := promhttp.Handler()
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := collector.CollectMetrics()
		if err != nil {
			slog.Error("Error collecting zettelkasten metrics", slog.Any("error", err))
			metrics.ExporterUp.Set(0)
		} else {
			metrics.ExporterUp.Set(1)
		}
		promHandler.ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	slog.Info("Starting HTTP server", slog.String("address", addr))
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("Error on HTTP server", slog.Any("error", err))
		os.Exit(1)
	}
}
