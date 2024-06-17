# Zettelkasten exporter

Prometheus exporter that collects statistics from your second brain. 

## TODO

- [X] Get zettelkasten from directory
- [X] Get total notes metric
- [X] Parse obsidian wiki links
- [X] Register total links
- [X] Register links per note
- [X] Expose prometheus metrics endpoint
- [X] Read config
- [X] Find all files recursivelly
- [X] Parse markdown links
- [X] Configurable ignore file patterns
- [X] Get zettelkasten from git url
- [X] Register metrics on InfluxDB
- [X] Make InfluxDB parameters configurable
- [X] Major refactor
- [X] Backfill data using git
- [ ] Support private repositories (Maybe with Github's PAT?)
- [ ] Handle InfluxDB async write errors (https://github.com/influxdata/influxdb-client-go?tab=readme-ov-file#reading-async-errors)
- [ ] Grafana dashboard
- [ ] Docker compose example
- [ ] Kubernetes example
- [ ] Document usage in README
- [ ] Github actions CI
- [ ] Build image and push to OCI registry
- [ ] Deploy on K8s
- [ ] Asynchronous git fetching

- [ ] Exclude links to non existing files
- [ ] Collect backlinks
- [ ] Collect word count
- [ ] Collect time to read

https://prometheus.io/docs/instrumenting/writing_exporters/
https://github.com/go-git/go-git/blob/master/_examples/pull/main.go
https://medium.com/tlvince/prometheus-backfilling-a92573eb712c
https://github.com/influxdata/helm-charts/tree/master/charts/influxdb2
https://grafana.com/docs/grafana/latest/getting-started/get-started-grafana-influxdb/
https://docs.influxdata.com/flux/v0/get-started/
