package collector

import (
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
)

type CollectorConfig struct {
	IgnorePatterns []string
}

// Collector represents a metrics collector.
type Collector struct {
	config  CollectorConfig
	storage storage.Storage
}

// NewCollector creates a new collector
func NewCollector(ignorePatterns []string, storage storage.Storage) Collector {
	return Collector{
		config: CollectorConfig{
			IgnorePatterns: ignorePatterns,
		},
		storage: storage,
	}
}

// CollectMetrics collects all metrics from a Zettelkasten rooted in `root` and writes them to the storage with a timestamp of `collectionTime`.
func (c *Collector) CollectMetrics(root fs.FS, collectionTime time.Time) error {
	slog.Debug("Collecting metrics", slog.Time("collection_time", collectionTime))
	start := time.Now()
	collected, err := c.collectMetrics(root)
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

// collectMetrics collects all metrics from a Zettelkasten rooted in `root`.
func (c *Collector) collectMetrics(root fs.FS) (metrics.Metrics, error) {
	noteMetrics := make(map[string]metrics.NoteMetrics)

	err := fs.WalkDir(root, ".", func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			slog.Error("Error on path. Will not enter it", slog.Any("error", err), slog.String("path", path))
			return nil
		}

		// Skip ignored files or directories
		if slices.Contains(c.config.IgnorePatterns, filepath.Base(path)) {
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
		defer f.Close()
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
		return metrics.Metrics{}, err
	}

	zettelkastenMetrics := aggregateMetrics(noteMetrics)
	return zettelkastenMetrics, nil
}

// aggregateMetrics aggregates all individual note metrics into metrics in the context of a full Zettelkasten.
func aggregateMetrics(noteMetrics map[string]metrics.NoteMetrics) metrics.Metrics {
	zettelkastenMetrics := metrics.Metrics{
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
