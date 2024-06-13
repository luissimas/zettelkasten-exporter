package metrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/luissimas/zettelkasten-exporter/internal/config"
)

var influxDB api.WriteAPI

func ConnectDatabase(cfg config.Config) {
	client := influxdb2.NewClient(cfg.InfluxDBURL, string(cfg.InfluxDBToken))
	influxDB = client.WriteAPI(cfg.InfluxDBOrg, cfg.InfluxDBBucket)
}

func RegisterNoteMetric(name string, linkCount int, timestamp time.Time) {
	point := influxdb2.NewPoint(
		name,
		map[string]string{"name": name},
		map[string]interface{}{"link_count": linkCount},
		timestamp,
	)
	influxDB.WritePoint(point)
}
