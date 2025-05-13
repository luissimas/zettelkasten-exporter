package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithVictoriaMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test")
	}

	ctx := context.Background()
	victoriametricsReq := testcontainers.ContainerRequest{
		Image:      "victoriametrics/victoria-metrics:v1.117.1",
		WaitingFor: wait.ForLog("started server at http://0.0.0.0:8428/"),
	}
	victoriametrics, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: victoriametricsReq,
		Started:          true,
	})
	defer testcontainers.CleanupContainer(t, victoriametrics)
	require.NoError(t, err)

	victoriametricsIP, err := victoriametrics.ContainerIP(ctx)
	require.NoError(t, err)

	victoriametricsPort := "8428"
	victoriametricsURL := fmt.Sprintf("http://%s:%s", victoriametricsIP, victoriametricsPort)
	require.NoError(t, err)
	env := map[string]string{
		"ZETTELKASTEN_DIRECTORY": "/sample",
		"VICTORIAMETRICS_URL":    victoriametricsURL,
		"COLLECTION_INTERVAL":    "5s",
	}
	_ = createExporterContainer(ctx, t, env)

	// TODO: assert that metrics were written propertly. Should we?
}

func TestWithInfluxDB(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test")
	}

	ctx := context.Background()
	influxDBReq := testcontainers.ContainerRequest{
		Image:      "influxdb:2",
		WaitingFor: wait.ForLog("Listening log_id="),
		Env: map[string]string{
			"DOCKER_INFLUXDB_INIT_MODE":        "setup",
			"DOCKER_INFLUXDB_INIT_USERNAME":    "admin",
			"DOCKER_INFLUXDB_INIT_PASSWORD":    "password",
			"DOCKER_INFLUXDB_INIT_ORG":         "default",
			"DOCKER_INFLUXDB_INIT_BUCKET":      "zettelkasten",
			"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN": "test-token",
		},
	}
	influxDB, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: influxDBReq,
		Started:          true,
	})
	defer testcontainers.CleanupContainer(t, influxDB)
	require.NoError(t, err)

	influxDBIP, err := influxDB.ContainerIP(ctx)
	require.NoError(t, err)

	influxDBPort := "8086"
	influxDBURL := fmt.Sprintf("http://%s:%s", influxDBIP, influxDBPort)
	require.NoError(t, err)
	env := map[string]string{
		"ZETTELKASTEN_DIRECTORY": "/sample",
		"INFLUXDB_TOKEN":         "test-token",
		"INFLUXDB_URL":           influxDBURL,
		"INFLUXDB_ORG":           "default",
		"INFLUXDB_BUCKET":        "zettelkasten",
		"COLLECTION_INTERVAL":    "5s",
	}
	_ = createExporterContainer(ctx, t, env)

	// TODO: assert that metrics were written propertly. Should we?
}

func createExporterContainer(ctx context.Context, t *testing.T, env map[string]string) testcontainers.Container {
	exporterReq := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Env:        env,
		WaitingFor: wait.ForLog("Starting metrics collection"),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "./examples/compose/sample",
				ContainerFilePath: "/sample",
			},
		},
	}
	exporter, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: exporterReq,
		Started:          true,
	})
	t.Cleanup(func() { testcontainers.CleanupContainer(t, exporter) })
	require.NoError(t, err)

	return exporter
}
