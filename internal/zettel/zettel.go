package zettel

import (
	"io/fs"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
)

// Zettel represents a zettelkasten
type Zettel interface {
	// Ensure makes sure that the zettelkasten is updated and operational.
	Ensure() error
	// GetRoot retrieves the root of the Zettelkasten directory structure.
	GetRoot() fs.FS
	// WalkHistory walks the history of the Zettelkasten, calling `walkFunc` on each point.
	WalkHistory(walkFunc func(time.Time) error) error
}

func NewZettel(cfg config.Config) Zettel {
	return NewGitZettel(cfg)
}
