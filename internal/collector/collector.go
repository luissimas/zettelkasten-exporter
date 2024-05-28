package collector

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

type Metrics struct {
	NoteCount int
	LinkCount int
	Notes     map[string]NoteMetrics
}

type CollectorConfig struct {
	Path string
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
			Path: absolute_path,
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
		metrics.LinkCount.WithLabelValues(name).Set(float64(metric.LinkCount))
	}

	elapsed := time.Since(started)
	metrics.CollectionDuration.Observe(float64(elapsed))
	slog.Info("Completed metrics collection", slog.Duration("duration", elapsed))
	return nil
}

func (c *Collector) collectMetrics() (Metrics, error) {
	// FIXME: filepath.Glob does not support double star expansion,
	// so this pattern is not searching recursivelly. We'll need to
	// walk the filesystem recursivelly.
	pattern := filepath.Join(c.config.Path, "**/*.md")
	files, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("Error getting files", slog.Any("error", err))
		return Metrics{}, err
	}

	noteCount := len(files)
	linkCount := 0
	notes := make(map[string]NoteMetrics)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			slog.Error("Error reading file", slog.Any("error", err))
			continue
		}
		metrics := CollectNoteMetrics(content)
		notes[file] = metrics
		linkCount += metrics.LinkCount
	}

	return Metrics{NoteCount: noteCount, LinkCount: linkCount, Notes: notes}, nil
}
