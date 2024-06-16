package storage

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

// The measurement name to be used for all metrics within the InfluxDB bucket.
const measurementName = "notes"

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

func (i InfluxDBStorage) WriteMetric(noteName string, metric metrics.NoteMetrics, timestamp time.Time) {
	point := influxdb2.NewPoint(
		measurementName,
		map[string]string{"name": noteName},
		map[string]interface{}{"link_count": metric.LinkCount},
		timestamp,
	)
	i.writeAPI.WritePoint(point)
}

func (i InfluxDBStorage) IsEmpty() bool {
	return false

}
