package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
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
	storage := storage.NewInfluxDBStorage(cfg.InfluxDBURL, cfg.InfluxDBOrg, cfg.InfluxDBBucket, cfg.InfluxDBToken)

	zettelkasten := zettel.NewZettel(cfg)
	collector := collector.NewCollector(zettelkasten.GetRoot(), cfg.IgnoreFiles, storage)
	err = zettelkasten.Ensure()
	if err != nil {
		slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
		os.Exit(1)
	}
	// TODO: check for empty bucket
	// slog.Info("Walking history")
	// start := time.Now()
	// err = zettelkasten.WalkHistory(collector.CollectMetrics)
	// if err != nil {
	// 	slog.Error("Error walking history", slog.Any("error", err))
	// 	os.Exit(1)
	// }
	// slog.Info("Collected historic metrics", slog.Duration("duration", time.Since(start)))
	for {
		err = zettelkasten.Ensure()
		if err != nil {
			slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
			os.Exit(1)
		}
		err = collector.CollectMetrics(time.Now())
		if err != nil {
			slog.Error("Error collecting metrics", slog.Any("error", err))
			os.Exit(1)
		}
		time.Sleep(cfg.CollectionInterval)
	}
}
