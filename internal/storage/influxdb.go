package storage

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// The measurement names to be used for metrics within the InfluxDB bucket.
const notesMeasurementName = "notes"
const totalMeasurementName = "total"

// InfluxDBStorage represents the implementation of a metric storage using InfluxDB.
type InfluxDBStorage struct {
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
}

// NewInfluxDBStorage creates a new `InfluxDBStorage`.
func NewInfluxDBStorage(url, org, bucket, token string) InfluxDBStorage {
	client := influxdb2.NewClient(url, string(token))
	writeAPI := client.WriteAPI(org, bucket)
	queryAPI := client.QueryAPI(org)
	return InfluxDBStorage{writeAPI: writeAPI, queryAPI: queryAPI}
}

// WriteMetric writes `metric` for `noteName` to the storage with `timestamp`.
func (i InfluxDBStorage) WriteMetrics(zettelkastenMetrics metrics.Metrics, timestamp time.Time) {
	// Aggregated metrics
	point := influxdb2.NewPoint(
		totalMeasurementName,
		map[string]string{},
		map[string]interface{}{
			"note_count": zettelkastenMetrics.NoteCount,
			"link_count": zettelkastenMetrics.LinkCount,
			"word_count": zettelkastenMetrics.WordCount,
		},
		timestamp,
	)
	i.writeAPI.WritePoint(point)

	// Individual note metrics
	for name, metric := range zettelkastenMetrics.Notes {
		point := influxdb2.NewPoint(
			notesMeasurementName,
			map[string]string{"name": name},
			map[string]interface{}{
				"link_count":     metric.LinkCount,
				"word_count":     metric.WordCount,
				"backlink_count": metric.BacklinkCount,
			},
			timestamp,
		)
		i.writeAPI.WritePoint(point)
	}
}
