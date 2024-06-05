package zettel

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
)

type LocalZettel struct {
	Config         config.Config
	LocalDirectory string
}

// NewLocalZettel creates a new LocalZettel
func NewLocalZettel(cfg config.Config) *LocalZettel {
	return &LocalZettel{Config: cfg, LocalDirectory: cfg.ZettelkastenDirectory}
}

// GetRoot retrieves the root of the local zettelkasten directory
func (l *LocalZettel) GetRoot() fs.FS {
	return os.DirFS(l.LocalDirectory)
}

// Ensure ensures that the local zettelkasten directory exists and is accessible.
func (l *LocalZettel) Ensure() error {
	if !filepath.IsAbs(l.LocalDirectory) {
		absolute_path, err := filepath.Abs(l.LocalDirectory)
		if err != nil {
			slog.Error("Error getting absolute path", slog.Any("error", err), slog.String("path", l.LocalDirectory))
			return err
		}
		l.LocalDirectory = absolute_path
	}
	_, err := os.Stat(l.LocalDirectory)
	if err != nil {
		slog.Error("Cannot stat zettelkasten directory", slog.Any("error", err), slog.String("path", l.LocalDirectory))
		return err
	}

	return nil
}
