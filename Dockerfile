FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /zettelkasten-exporter

FROM alpine:3 AS release-stage

WORKDIR /

COPY --from=build-stage /zettelkasten-exporter /zettelkasten-exporter

RUN apk add git

ENTRYPOINT ["/zettelkasten-exporter"]
