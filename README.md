# Zettelkasten exporter

> [!WARNING]
> This is still a work in progress, expect breaking changes

An agent that collects metrics from an zettelkasten and stores them into an InfluxDB bucket.

![](./docs/assets/dashboard.png)

## Features

- Collect metrics from a local directory or a git repository
- Backfill historical metrics
- Parses both markdown and wiki links
- Authenticate in private git repositories using personal access tokens
- Grafana dashboards included

## References

https://prometheus.io/docs/instrumenting/writing_exporters/

https://github.com/go-git/go-git/blob/master/_examples/pull/main.go

https://medium.com/tlvince/prometheus-backfilling-a92573eb712c

https://github.com/influxdata/helm-charts/tree/master/charts/influxdb2

https://grafana.com/docs/grafana/latest/getting-started/get-started-grafana-influxdb/

https://docs.influxdata.com/flux/v0/get-started/
