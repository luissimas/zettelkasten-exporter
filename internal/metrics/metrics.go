package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CollectionSuccessful = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "zettelkasten_collection_successful",
		Help: "Whether the metrics collections were successful",
	})
	TotalNoteCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "zettelkasten_total_note_count",
		Help: "The total count of notes in the zettelkasten",
	})
	LinkCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "zettelkasten_link_count",
		Help: "The count of links in the zettelkasten",
	}, []string{"name"})
	CollectionDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "zettelkasten_collection_duration",
		Help:    "The duration of the metrics collection",
		Buckets: prometheus.LinearBuckets(1, 1000, 10),
	})
)
