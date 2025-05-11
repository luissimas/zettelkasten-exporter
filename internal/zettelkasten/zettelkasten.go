package zettelkasten

import (
	"io/fs"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
)

// Zettelkasten represents a Zettelkasten.
type Zettelkasten interface {
	// Ensure makes sure that the zettelkasten is updated and operational.
	Ensure() error
	// GetRoot retrieves the root of the Zettelkasten directory structure.
	GetRoot() fs.FS
	// WalkHistory walks the history of the Zettelkasten, calling `walkFunc` on each point in the history.
	WalkHistory(walkFunc WalkFunc) error
}

// WalkFunc is the type of function called by `Zettelkasten.WalkHistory` to
// process all points in the history of the zettelkasten.
type WalkFunc func(root fs.FS, timestamp time.Time) error

// NewZettelkasten creates a new Zettelkasten from the given config.
func NewZettelkasten(cfg config.Config) Zettelkasten {
	if cfg.ZettelkastenGitURL != "" {
		return NewGitZettelkasten(cfg.ZettelkastenGitURL, cfg.ZettelkastenGitBranch, cfg.ZettelkastenGitToken)
	}

	return NewLocalZettelkasten(cfg.ZettelkastenDirectory)
}
