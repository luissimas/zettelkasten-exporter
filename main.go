package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/exporter"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
	"github.com/luissimas/zettelkasten-exporter/internal/zettelkasten"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Debug("Loaded config", slog.Any("config", cfg))
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	metricsStorage, err := storage.NewStorage(cfg)
	if err != nil {
		slog.Error("Error creating storage", slog.Any("error", err))
		os.Exit(1)
	}

	zet := zettelkasten.NewZettelkasten(cfg)
	exporter := exporter.NewExporter(cfg, zet, metricsStorage)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := exporter.Start(ctx); err != nil {
		slog.Error("Error on exporter", slog.Any("error", err))
		os.Exit(1)
	}
}
