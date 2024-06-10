package metrics

import (
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var influxDB api.WriteAPI

func ConnectDatabase() {
	// TODO: get values from config
	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)
	org := "default"
	bucket := "zettelkasten"
	influxDB = client.WriteAPI(org, bucket)
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
