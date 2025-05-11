package storage

import (
	"errors"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// Storage represents a storage for metrics.
type Storage interface {
	// WriteMetric writes the `zettelkastenMetrics` to the storage.
	WriteMetrics(zettelkastenMetrics metrics.ZettelkastenMetrics, timestamp time.Time) error
}

// NewStorage creates a new Storage from the given config.
func NewStorage(cfg config.Config) (Storage, error) {
	if cfg.VictoriaMetricsURL != "" {
		return NewVictoriaMetricsStorage(cfg.VictoriaMetricsURL), nil
	}

	if cfg.InfluxDBURL != "" {
		return NewInfluxDBStorage(cfg.InfluxDBURL, cfg.InfluxDBOrg, cfg.InfluxDBBucket, cfg.InfluxDBToken), nil
	}

	return nil, errors.New("invalid storage config")
}
