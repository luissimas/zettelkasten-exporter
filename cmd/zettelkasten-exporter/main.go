package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/luissimas/zettelkasten-exporter/internal/zettel"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Debug("Loaded config", slog.Any("config", cfg))
	zettelkasten := zettel.NewZettel(cfg)
	err = zettelkasten.Ensure()
	if err != nil {
		slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
		os.Exit(1)
	}

	collector := collector.NewCollector(zettelkasten.GetRoot(), cfg.IgnoreFiles)
	promHandler := promhttp.Handler()
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Starting metrics collection")
		started := time.Now()
		err := zettelkasten.Ensure()
		if err != nil {
			slog.Error("Error ensuring zettelkasten", slog.Any("error", err))
			metrics.ExporterUp.Set(0)
		} else {
			err = collector.CollectMetrics()
			if err != nil {
				slog.Error("Error collecting zettelkasten metrics", slog.Any("error", err))
				metrics.ExporterUp.Set(0)
			} else {
				metrics.ExporterUp.Set(1)
			}
		}

		elapsed := time.Since(started)
		metrics.CollectionDuration.Observe(float64(elapsed))
		slog.Info("Completed metrics collection", slog.Duration("duration", elapsed))

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
