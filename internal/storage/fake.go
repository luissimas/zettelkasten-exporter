package storage

import (
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// FakeStorage represents a fake implementation of storage to be used in tests.
type FakeStorage struct{}

// FakeStorage creates a new `FakeStorage`.
func NewFakeStorage() FakeStorage {
	return FakeStorage{}
}

func (f FakeStorage) WriteMetrics(zettelkastenMetrics metrics.Metrics, timestamp time.Time) error {
	return nil
}

func (f FakeStorage) IsEmpty() bool {
	return false
}
