FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /zettelkasten-exporter ./cmd/zettelkasten-exporter/main.go

FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /zettelkasten-exporter /zettelkasten-exporter

USER nonroot:nonroot

ENTRYPOINT ["/zettelkasten-exporter"]
