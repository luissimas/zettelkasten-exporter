package storage

import (
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// FakeStorage represents a fake implementation of storage to be used in tests.
type FakeStorage struct {
	Metrics []metrics.ZettelkastenMetrics
}

// FakeStorage creates a new `FakeStorage`.
func NewFakeStorage() FakeStorage {
	return FakeStorage{}
}

func (f *FakeStorage) WriteMetrics(zettelkastenMetrics metrics.ZettelkastenMetrics, timestamp time.Time) error {
	f.Metrics = append(f.Metrics, zettelkastenMetrics)
	return nil
}
