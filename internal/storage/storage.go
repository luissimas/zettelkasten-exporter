package storage

import (
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// Storage represents a storage for metrics.
type Storage interface {
	// WriteMetric writes the `zettelkastenMetrics` to the storage.
	WriteMetrics(zettelkastenMetrics metrics.Metrics, timestamp time.Time) error
}
