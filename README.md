[![CI](https://github.com/luissimas/zettelkasten-exporter/actions/workflows/check.yaml/badge.svg)](https://github.com/luissimas/zettelkasten-exporter/actions/workflows/check.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/luissimas/zettelkasten-exporter)](https://goreportcard.com/report/github.com/luissimas/zettelkasten-exporter)
[![Codecov](https://codecov.io/github/luissimas/zettelkasten-exporter/coverage.svg?branch=main)](https://codecov.io/gh/luissimas/zettelkasten-exporter)

# Zettelkasten exporter

An agent that collects metrics from an zettelkasten and lets you visualize them in Grafana.

![](./docs/assets/dashboard.png)

> [!NOTE]
> This is still a work in progress. If you find any bugs or have any questions, feel free to [open an issue](https://github.com/luissimas/zettelkasten-exporter/issues/new/choose).

## Features

- Collect metrics from a local directory or a git repository
- Backfill historical metrics using git
- Parses both markdown and wiki links
- Authenticate in private git repositories using personal access tokens
- Grafana dashboards included
- Support for both InfluxDB and VictoriaMetrics as storage backends

## Usage

The exporter is distributed as both a binary and a Docker image. The currently documented ways of deploy are via docker compose and Kubernetes. For details about the setup, check out the examples:

- [Docker compose example](./examples/compose)
- [Kubernetes example](./examples/kubernetes)

Note that for a complete solution, it will be necessary to setup Grafana and either InfluxDB or VictoriaMetrics as data sources. For more information about setting them up, refer to their documentation. Here are some links that might be useful:

- https://grafana.com/docs/grafana/latest/getting-started/get-started-grafana-influxdb/
- https://docs.influxdata.com/influxdb/v2/get-started/setup/
- https://docs.victoriametrics.com/

In the `dashboards` folder there are two Grafana dashboards provided: one for InfluxDB and another for VictoriaMetrics. You'll have to import the appropriate dashboard for your storage backend.

The provided InfluxDB dashboard uses `Flux` as the query language, so make sure to set the "Query language" option to "Flux" when creating the InfluxDB data source in Grafana.

For both storage backends, make sure to configure the data retention period according to your needs.

## Configuration

All configuration is supplied via environment variables. You should supply at least the zettelkasten source via the `ZETTELKASTEN_DIRECTORY` or `ZETTELKASTEN_GIT_URL` variables and the storage backend via the `VICTORIAMETRICS_URL` or `INFLUXDB_*` variables.

| Name                       | Description                                                          | Default                        | Required |
| -------------------------- | -------------------------------------------------------------------- | ------------------------------ | -------- |
| VICTORIAMETRICS_URL        | The VictoriaMetrics URL                                              |                                | No       |
| INFLUXDB_URL               | The InfluxDB URL                                                     |                                | No       |
| INFLUXDB_TOKEN             | The InfluxDB token to authenticate in the bucket                     |                                | No       |
| INFLUXDB_ORG               | The InfluxDB org containing the bucket                               |                                | No       |
| INFLUXDB_BUCKET            | The InfluxDB bucket to register metrics                              |                                | No       |
| ZETTELKASTEN_DIRECTORY     | The local directory containing the zettelkasten                      |                                | No       |
| ZETTELKASTEN_GIT_URL       | The URL for the git repository containing the zettelkasten           |                                | No       |
| ZETTELKASTEN_GIT_TOKEN     | The access token to authenticate with private repositories           |                                | No       |
| ZETTELKASTEN_GIT_BRANCH    | The branch to use for git repositories                               | main                           | No       |
| COLLECTION_INTERVAL        | Time to wait between metric collections                              | 5m                             | No       |
| COLLECT_HISTORICAL_METRICS | Wether to collect historical metrics at startup                      | true                           | No       |
| IGNORE_FILES               | Comma separated list of files that will be ignored in the collection | .git,obsidian,.trash,README.md | No       |
| LOG_LEVEL                  | The minimum log level                                                | INFO                           | No       |

## Metrics

The exporter collects metrics by parsing the contents of the markdown files present in the Zettelkasten. Currently the exporter stores metrics for individual notes and also aggregated metrics describing the entire Zettelkasten. The combination of raw and pre processed metrics allows for both flexibility and efficiency when querying the data, at the cost of a slightly higher storage usage. When using the InfluxDB storage, the two sets of metrics are stored in the same InfluxDB bucket under different [measurement names](https://docs.influxdata.com/influxdb/cloud/reference/key-concepts/data-elements/#measurement). When using the VictoriaMetrics storage, each metric is stored under a different name.

The following table describes all metrics collected by the exporter and their respective measurement names:

| InfluxDB measurement | InfluxDB name  | VictoriaMetrics name | Description                             |
|----------------------|----------------|----------------------|-----------------------------------------|
| notes                | link_count     | notes_link_count     | Number of links in the note             |
| notes                | word_count     | notes_word_count     | Number of words in the note             |
| notes                | backlink_count | notes_backlink_count | Number of links that reference the note |
| total                | note_count     | total_note_count     | Number of notes in the Zettelkasten     |
| total                | link_count     | total_link_count     | Number of links in the Zettelkasten     |
| total                | word_count     | total_word_count     | Number of words in the Zettelkasten     |

## Roadmap

These are some features that I'd like to include in the future.

- Support Prometheus remote write as a storage backend

## References

https://prometheus.io/docs/instrumenting/writing_exporters/

https://medium.com/tlvince/prometheus-backfilling-a92573eb712c

https://github.com/influxdata/helm-charts/tree/master/charts/influxdb2

https://grafana.com/docs/grafana/latest/getting-started/get-started-grafana-influxdb/

https://docs.influxdata.com/flux/v0/get-started/

https://github.com/onedr0p/exportar

https://github.com/mischavandenburg/zettelkasten-tracker
