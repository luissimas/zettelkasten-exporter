package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/luissimas/zettelkasten-exporter/internal/collector"
)

func main() {
	if len(os.Args) != 2 {
		slog.Error("Wrong arguments")
		os.Exit(1)
	}

	dir := os.Args[1]
	absolute_path, err := filepath.Abs(dir)
	if err != nil {
		slog.Error("Error getting absolute file path", slog.Any("error", err))
		os.Exit(1)
	}

	metrics, err := collector.CollectMetrics(collector.CollectorConfig{Path: absolute_path})
	if err != nil {
		slog.Error("Error collecting note metrics", slog.Any("error", err))
		os.Exit(1)
	}

	fmt.Printf("There are %d markdown files in %s\n", metrics.NoteCount, absolute_path)
	fmt.Printf("There is a total of %d links in%s\n", metrics.LinkCount, absolute_path)
	for name, metrics := range metrics.Notes {
		fmt.Printf("\t%s: %v\n", name, metrics)
	}
}
