package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
	"github.com/luissimas/zettelkasten-exporter/internal/zettelkasten"
)

func main() {
	// Setup
	cfg, err := config.LoadConfig()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Debug("Loaded config", slog.Any("config", cfg))
	storage := storage.NewInfluxDBStorage(cfg.InfluxDBURL, cfg.InfluxDBOrg, cfg.InfluxDBBucket, cfg.InfluxDBToken)
	collector := collector.NewCollector(cfg.IgnoreFiles, storage)
	var zet zettelkasten.Zettelkasten
	if cfg.ZettelkastenGitURL != "" {
		zet = zettelkasten.NewGitZettelkasten(cfg.ZettelkastenGitURL, cfg.ZettelkastenGitBranch, cfg.ZettelkastenGitToken)
	} else {
		zet = zettelkasten.NewLocalZettelkasten(cfg.ZettelkastenDirectory)
	}

	// Collect historical data
	if cfg.CollectHistoricalMetrics {
		slog.Info("Collecting historical metrics")
		start := time.Now()
		err = zet.Ensure()
		if err != nil {
			slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Walking zettelkasten history")
		err = zet.WalkHistory(collector.CollectMetrics)
		if err != nil {
			slog.Error("Error walking history", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Collected historical metrics", slog.Duration("duration", time.Since(start)))
	}

	// Periodic collection loop
	for {
		slog.Info("Starting metrics collection")
		start := time.Now()
		err = zet.Ensure()
		if err != nil {
			slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
			os.Exit(1)
		}
		root := zet.GetRoot()
		err = collector.CollectMetrics(root, time.Now())
		if err != nil {
			slog.Error("Error collecting metrics", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Collected metrics", slog.Duration("duration", time.Since(start)))
		time.Sleep(cfg.CollectionInterval)
	}
}
