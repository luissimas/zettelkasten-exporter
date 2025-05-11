package exporter

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
	"github.com/luissimas/zettelkasten-exporter/internal/zettelkasten"
)

// Exporter represents a Zettelkasten exporter.
type Exporter struct {
	config       config.Config
	storage      storage.Storage
	zettelkasten zettelkasten.Zettelkasten
	ticker       *time.Ticker
}

// NewExporter creates a new exporter.
func NewExporter(cfg config.Config, zettelkasten zettelkasten.Zettelkasten, storage storage.Storage) Exporter {
	return Exporter{
		config:       cfg,
		storage:      storage,
		zettelkasten: zettelkasten,
		ticker:       time.NewTicker(cfg.CollectionInterval),
	}
}

// Start starts the exporter loop.
func (e *Exporter) Start(ctx context.Context) {
	// Collect historical data
	if e.config.CollectHistoricalMetrics {
		slog.Info("Collecting historical metrics")
		start := time.Now()
		err := e.zettelkasten.Ensure()
		if err != nil {
			slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
		}

		slog.Info("Walking zettelkasten history")
		err = e.zettelkasten.WalkHistory(e.collectMetrics)
		if err != nil {
			slog.Error("Error walking history", slog.Any("error", err))
		}

		slog.Info("Collected historical metrics", slog.Duration("duration", time.Since(start)))
	}

	for {
		select {
		case t := <-e.ticker.C:
			slog.Info("Starting metrics collection")
			err := e.zettelkasten.Ensure()
			if err != nil {
				slog.Error("Error ensuring that zettelkasten is ready", slog.Any("error", err))
			}

			err = e.collectMetrics(e.zettelkasten.GetRoot(), t)
			if err != nil {
				slog.Error("Error collecting metrics", slog.Any("error", err))
			}

			slog.Info("Collected metrics", slog.Duration("duration", time.Since(t)), slog.Time("next_run", time.Now().Add(e.config.CollectionInterval)))
		case <-ctx.Done():
			slog.Info("Stopping metrics collection")
			return
		}
	}
}

// collectMetrics collects all metrics from a Zettelkasten rooted in `root` and writes them to the storage with a timestamp of `collectionTime`.
func (c *Exporter) collectMetrics(root fs.FS, collectionTime time.Time) error {
	slog.Debug("Collecting metrics", slog.Time("collection_time", collectionTime))
	start := time.Now()
	collected, err := c.scrapeMetrics(root)
	if err != nil {
		return err
	}

	err = c.storage.WriteMetrics(collected, collectionTime)
	if err != nil {
		return err
	}

	slog.Debug("Collected metrics", slog.Duration("duration", time.Since(start)))

	return nil
}

// scrapeMetrics collects all metrics from a Zettelkasten rooted in `root`.
func (c *Exporter) scrapeMetrics(root fs.FS) (metrics.ZettelkastenMetrics, error) {
	noteMetrics := make(map[string]metrics.NoteMetrics)

	err := fs.WalkDir(root, ".", func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("Error on path. Will not enter it", slog.Any("error", err), slog.String("path", path))
			return nil
		}

		// Skip ignored files or directories
		if slices.Contains(c.config.IgnoreFiles, filepath.Base(path)) {
			if dir.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip other directories and non markdown files
		if dir.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		f, err := root.Open(path)
		if err != nil {
			slog.Error("Error opening file", slog.Any("error", err), slog.String("path", path))
			return nil
		}

		defer func() {
			err := f.Close()
			if err != nil {
				slog.Warn("Error closing file", slog.Any("error", err), slog.String("path", path))
			}
		}()
		content, err := io.ReadAll(f)
		if err != nil {
			slog.Error("Error reading file", slog.Any("error", err), slog.String("path", path))
			return nil
		}

		noteMetrics[nameFromFilename(path)] = CollectNoteMetrics(content)
		slog.Debug("collected metrics from file", slog.String("path", path), slog.Any("d", dir), slog.Any("err", err))

		return nil
	})

	if err != nil {
		slog.Error("Error getting files", slog.Any("error", err))
		return metrics.ZettelkastenMetrics{}, err
	}

	zettelkastenMetrics := aggregateMetrics(noteMetrics)
	return zettelkastenMetrics, nil
}

// aggregateMetrics aggregates all individual note metrics into metrics in the context of a full Zettelkasten.
func aggregateMetrics(noteMetrics map[string]metrics.NoteMetrics) metrics.ZettelkastenMetrics {
	zettelkastenMetrics := metrics.ZettelkastenMetrics{
		NoteCount: 0,
		LinkCount: 0,
		WordCount: 0,
		Notes:     make(map[string]metrics.NoteMetrics),
	}

	for name, metric := range noteMetrics {
		// Aggregate totals
		zettelkastenMetrics.NoteCount += 1
		zettelkastenMetrics.LinkCount += metric.LinkCount
		zettelkastenMetrics.WordCount += metric.WordCount
		// Collect backlinks
		for _, n := range noteMetrics {
			metric.BacklinkCount += n.Links[name]
		}
		zettelkastenMetrics.Notes[name] = metric
	}

	return zettelkastenMetrics
}
