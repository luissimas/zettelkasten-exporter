package storage

import (
	"context"
	"log/slog"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// The measurement names to be used for metrics within the InfluxDB bucket.
const notesMeasurementName = "notes"
const totalMeasurementName = "total"

// InfluxDBStorage represents the implementation of a metric storage using InfluxDB.
type InfluxDBStorage struct {
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
}

// NewInfluxDBStorage creates a new `InfluxDBStorage`.
func NewInfluxDBStorage(url, org, bucket, token string) InfluxDBStorage {
	client := influxdb2.NewClient(url, string(token))
	writeAPI := client.WriteAPIBlocking(org, bucket)
	queryAPI := client.QueryAPI(org)
	return InfluxDBStorage{writeAPI: writeAPI, queryAPI: queryAPI}
}

// WriteMetric writes `metric` for `noteName` to the storage with `timestamp`.
func (i InfluxDBStorage) WriteMetrics(zettelkastenMetrics metrics.ZettelkastenMetrics, timestamp time.Time) error {
	points := createInfluxDBPoints(zettelkastenMetrics, timestamp)
	slog.Debug("Writing metrics to InfluxDB", slog.Any("points", points))
	err := i.writeAPI.WritePoint(context.Background(), points...)
	if err != nil {
		slog.Error("Error writing points to InfluxDB storage", slog.Any("error", err))
	}
	return err
}

// createInfluxDBPoints creates a slice of InfluxDB measurement points from `zettelkastenMetrics` with the given `timestamp`.
func createInfluxDBPoints(zettelkastenMetrics metrics.ZettelkastenMetrics, timestamp time.Time) []*write.Point {
	points := make([]*write.Point, 0, len(zettelkastenMetrics.Notes)+1)
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
	points = append(points, point)

	// Individual note metrics
	for name, metric := range zettelkastenMetrics.Notes {
		point = influxdb2.NewPoint(
			notesMeasurementName,
			map[string]string{"name": name},
			map[string]interface{}{
				"link_count":     metric.LinkCount,
				"word_count":     metric.WordCount,
				"backlink_count": metric.BacklinkCount,
			},
			timestamp,
		)
		points = append(points, point)
	}
	return points
}
