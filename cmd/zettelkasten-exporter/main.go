package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/luissimas/zettelkasten-exporter/internal/zettel"
)

func main() {
	cfg, err := config.LoadConfig()
	slog.SetLogLoggerLevel(cfg.LogLevel)
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Debug("Loaded config", slog.Any("config", cfg))
	metrics.ConnectDatabase()
	zettelkasten := zettel.NewZettel(cfg)
	err = zettelkasten.Ensure()
	if err != nil {
		slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
		os.Exit(1)
	}

	collector := collector.NewCollector(zettelkasten.GetRoot(), cfg.IgnoreFiles)

	for {
		err := collector.CollectMetrics()
		if err != nil {
			slog.Error("Error collecting metrics", slog.Any("error", err))
			os.Exit(1)
		}
		time.Sleep(10 * time.Second)
	}
}
