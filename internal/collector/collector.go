package collector

import (
	"io"
	"io/fs"
	"log/slog"
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
	FileSystem     fs.FS
	IgnorePatterns []string
}

type Collector struct {
	config CollectorConfig
}

func NewCollector(fileSystem fs.FS, ignorePatterns []string) Collector {
	return Collector{
		config: CollectorConfig{
			FileSystem:     fileSystem,
			IgnorePatterns: ignorePatterns,
		},
	}
}

func (c *Collector) CollectMetrics() error {
	slog.Info("Collecting metrics")
	start := time.Now()
	collected, err := c.collectMetrics()
	if err != nil {
		return err
	}

	collectionTime := time.Now()
	for name, metric := range collected.Notes {
		metrics.RegisterNoteMetric(name, metric.LinkCount, collectionTime)
	}
	slog.Info("Collected metrics", slog.Duration("duration", time.Since(start)))

	return nil
}

func (c *Collector) collectMetrics() (Metrics, error) {
	noteCount := 0
	linkCount := 0
	notes := make(map[string]NoteMetrics)

	err := fs.WalkDir(c.config.FileSystem, ".", func(path string, dir fs.DirEntry, err error) error {
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

		f, err := c.config.FileSystem.Open(path)
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
		return Metrics{}, err
	}

	return Metrics{NoteCount: noteCount, LinkCount: linkCount, Notes: notes}, nil
}
