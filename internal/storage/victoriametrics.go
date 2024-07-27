package storage

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	lp "github.com/influxdata/line-protocol"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
)

type VictoriaMetricsStorage struct {
	writeUrl string
}

func NewVictoriaMetricsStorage(url string) VictoriaMetricsStorage {
	return VictoriaMetricsStorage{writeUrl: fmt.Sprintf("%s/api/v2/write", url)}
}

func (v VictoriaMetricsStorage) WriteMetrics(zettelkastenMetrics metrics.Metrics, timestamp time.Time) {
	// NOTE: we encode the metrics in the InfluxDB line protocol and write them to the VictoriaMetrics write endpoint.
	// Reference: https://docs.victoriametrics.com/#how-to-send-data-from-influxdb-compatible-agents-such-as-telegraf
	points := createInfluxDBPoints(zettelkastenMetrics, timestamp)
	content, err := encodePoints(points)
	if err != nil {
		slog.Error("Error encoding points into line procotol", slog.Any("error", err))
	}
	slog.Info("Writing metrics to endpoint", slog.String("content", string(content)))
	_, err = http.Post(v.writeUrl, "application/x-www-form-urlencoded", bytes.NewBuffer(content))
	if err != nil {
		slog.Error("Error sending POST request to endpoint", slog.Any("error", err), slog.String("url", v.writeUrl))
	}
}

// encodePoints encodes the given `points` into InfluxDB's line protocol.
func encodePoints(points []*write.Point) ([]byte, error) {
	var buffer bytes.Buffer
	e := lp.NewEncoder(&buffer)
	e.SetFieldTypeSupport(lp.UintSupport)
	e.FailOnFieldErr(true)
	e.SetPrecision(time.Millisecond)
	slog.Info("Endcoding points", slog.Any("points", points))
	for _, point := range points {
		_, err := e.Encode(point)
		if err != nil {
			return make([]byte, 0), err
		}
	}
	return buffer.Bytes(), nil
}
