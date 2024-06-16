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

func (f FakeStorage) WriteMetric(noteName string, metric metrics.NoteMetrics, timestamp time.Time) {

}

func (f FakeStorage) IsEmpty() bool {
	return false
}
