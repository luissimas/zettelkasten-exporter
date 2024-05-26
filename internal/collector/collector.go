package collector

import (
	"log/slog"
	"os"
	"path/filepath"
)

type Metrics struct {
	NoteCount int
	LinkCount int
	Notes     map[string]NoteMetrics
}

type CollectorConfig struct {
	Path string
}

func CollectMetrics(config CollectorConfig) (Metrics, error) {
	pattern := filepath.Join(config.Path, "**/*.md")
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
