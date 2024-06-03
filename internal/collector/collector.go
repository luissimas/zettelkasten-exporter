package collector

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

type Metrics struct {
	NoteCount int
	LinkCount int
	Notes     map[string]NoteMetrics
}

type CollectorConfig struct {
	Path           string
	IgnorePatterns []string
}

type Collector struct {
	config CollectorConfig
}

func NewCollector(path string) (Collector, error) {
	absolute_path, err := filepath.Abs(path)
	if err != nil {
		return Collector{}, err
	}

	return Collector{
		config: CollectorConfig{
			Path:           absolute_path,
			IgnorePatterns: []string{".obsidian"},
		},
	}, nil
}

func (c *Collector) CollectMetrics() error {
	started := time.Now()
	slog.Info("Starting metrics collection")

	collected, err := c.collectMetrics()
	if err != nil {
		return err
	}

	metrics.TotalNoteCount.Set(float64(collected.NoteCount))
	for name, metric := range collected.Notes {
		metrics.LinkCount.WithLabelValues(name).Set(float64(len(metric.Links)))
	}

	elapsed := time.Since(started)
	metrics.CollectionDuration.Observe(float64(elapsed))
	slog.Info("Completed metrics collection", slog.Duration("duration", elapsed))
	return nil
}

func (c *Collector) collectMetrics() (Metrics, error) {
	_, err := os.Stat(c.config.Path)
	if err != nil && os.IsNotExist(err) {
		return Metrics{}, err
	}

	noteCount := 0
	linkCount := 0
	notes := make(map[string]NoteMetrics)

	filepath.WalkDir(c.config.Path, func(path string, d fs.DirEntry, err error) error {
		// Skip all files in ignored directories
		if slices.Contains(c.config.IgnorePatterns, filepath.Base(path)) {
			return filepath.SkipDir
		}
		// Skip other directories and non markdown files
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			slog.Error("Error reading file", slog.Any("error", err))
			return nil
		}
		metrics := CollectNoteMetrics(content)
		notes[path] = metrics
		linkCount += len(metrics.Links)
		noteCount += 1

		slog.Info("collected metrics from file", slog.String("path", path), slog.Any("d", d), slog.Any("err", err))

		return nil
	})

	if err != nil {
		slog.Error("Error getting files", slog.Any("error", err))
		return Metrics{}, err
	}

	return Metrics{NoteCount: noteCount, LinkCount: linkCount, Notes: notes}, nil
}
