package storage

import (
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// Storage represents a storage for metrics.
type Storage interface {
	// WriteMetric writes the note metric to the storage.
	WriteMetric(noteName string, metric metrics.NoteMetrics, timestamp time.Time)
}
