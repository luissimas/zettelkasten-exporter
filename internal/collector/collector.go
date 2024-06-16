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

type Collector struct {
	config  CollectorConfig
	storage storage.Storage
}

func NewCollector(ignorePatterns []string, storage storage.Storage) Collector {
	return Collector{
		config: CollectorConfig{
			IgnorePatterns: ignorePatterns,
		},
		storage: storage,
	}
}

func (c *Collector) CollectMetrics(root fs.FS, collectionTime time.Time) error {
	slog.Info("Collecting metrics")
	start := time.Now()
	collected, err := c.collectMetrics(root)
	if err != nil {
		return err
	}

	for name, metric := range collected.Notes {
		c.storage.WriteMetric(name, metric, collectionTime)
	}
	slog.Info("Collected metrics", slog.Duration("duration", time.Since(start)))

	return nil
}

func (c *Collector) collectMetrics(root fs.FS) (metrics.Metrics, error) {
	noteCount := 0
	linkCount := 0
	notes := make(map[string]metrics.NoteMetrics)

	err := fs.WalkDir(root, ".", func(path string, dir fs.DirEntry, err error) error {
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
		content, err := io.ReadAll(f)
		if err != nil {
			slog.Error("Error reading file", slog.Any("error", err))
			return nil
		}
		metrics := CollectNoteMetrics(content)
		notes[path] = metrics
		linkCount += metrics.LinkCount
		noteCount += 1

		slog.Debug("collected metrics from file", slog.String("path", path), slog.Any("d", dir), slog.Any("err", err))

		return nil
	})

	if err != nil {
		slog.Error("Error getting files", slog.Any("error", err))
		return metrics.Metrics{}, err
	}

	return metrics.Metrics{NoteCount: noteCount, LinkCount: linkCount, Notes: notes}, nil
}
