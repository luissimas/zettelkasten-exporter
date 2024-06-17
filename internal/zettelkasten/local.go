package zettelkasten

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// LocalZettelkasten represents a zettelkasten locally with no version control.
type LocalZettelkasten struct {
	rootPath string
}

// NewLocalZettelkasten creates a new LocalZettelkasten.
func NewLocalZettelkasten(rootPath string) LocalZettelkasten {
	return LocalZettelkasten{rootPath: rootPath}
}

// GetRoot retrieves the root of the local zettelkasten directory
func (l LocalZettelkasten) GetRoot() fs.FS {
	return os.DirFS(l.rootPath)
}

// Ensure ensures that the local zettelkasten directory exists and is accessible.
func (l LocalZettelkasten) Ensure() error {
	if !filepath.IsAbs(l.rootPath) {
		absolute_path, err := filepath.Abs(l.rootPath)
		if err != nil {
			slog.Error("Error getting absolute path", slog.Any("error", err), slog.String("path", l.rootPath))
			return err
		}
		l.rootPath = absolute_path
	}
	_, err := os.Stat(l.rootPath)
	if err != nil {
		slog.Error("Cannot stat zettelkasten directory", slog.Any("error", err), slog.String("path", l.rootPath))
		return err
	}

	return nil
}

// WalkHistory calls `walkFunc` for each point in the zettelkasten history.
func (l LocalZettelkasten) WalkHistory(walkFunc WalkFunc) error {
	return walkFunc(l.GetRoot(), time.Now())
}
