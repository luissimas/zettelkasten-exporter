BINARY_NAME=zettelkasten-exporter

.PHONY: all clean format test build run

all: format vet build

clean:
	go clean
	rm bin/$(BINARY_NAME)

format:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build:
	go build -o bin/$(BINARY_NAME) ./cmd/zettelkasten-exporter/main.go

run: build
	./bin/$(BINARY_NAME) ~/Documents/vault/
