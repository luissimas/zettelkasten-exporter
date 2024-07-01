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

	// Write metrics to storage
	for name, metric := range collected.Notes {
		c.storage.WriteMetric(name, metric, collectionTime)
	}
	slog.Debug("Collected metrics", slog.Duration("duration", time.Since(start)))

	return nil
}

// collectMetrics collects all metrics from a Zettelkasten rooted in `root`.
func (c *Collector) collectMetrics(root fs.FS) (metrics.Metrics, error) {
	var noteCount uint
	var linkCount uint
	var wordCount uint
	notes := make(map[string]metrics.NoteMetrics)

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
		metrics := CollectNoteMetrics(content)
		notes[path] = metrics
		linkCount += metrics.LinkCount
		wordCount += metrics.WordCount
		noteCount += 1

		slog.Debug("collected metrics from file", slog.String("path", path), slog.Any("d", dir), slog.Any("err", err))

		return nil
	})

	if err != nil {
		slog.Error("Error getting files", slog.Any("error", err))
		return metrics.Metrics{}, err
	}

	return metrics.Metrics{NoteCount: noteCount, LinkCount: linkCount, WordCount: wordCount, Notes: notes}, nil
}
