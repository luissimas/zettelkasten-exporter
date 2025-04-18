BINARY_NAME=zettelkasten-exporter

.PHONY: all clean format vet lint test build run

all: format lint vet build test

clean:
	go clean
	rm bin/$(BINARY_NAME)

format:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./...

build:
	go build -o bin/$(BINARY_NAME)

run: build
	LOG_LEVEL=INFO \
	ZETTELKASTEN_GIT_URL=<GIT_URL> \
	ZETTELKASTEN_GIT_BRANCH=master \
	ZETTELKASTEN_GIT_TOKEN=<GIT_TOKEN> \
	COLLECTION_INTERVAL=10s \
	INFLUXDB_TOKEN=<INFLUXDB_TOKEN> \
	INFLUXDB_URL=http://localhost:8086 \
	INFLUXDB_ORG=default \
	INFLUXDB_BUCKET=zettelkasten \
	./bin/$(BINARY_NAME)

docker:
	docker build . -t zettelkasten-exporter:latest
