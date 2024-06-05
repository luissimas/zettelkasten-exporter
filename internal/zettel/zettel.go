package zettel

import (
	"io/fs"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
)

// Zettel represents a zettelkasten
type Zettel interface {
	// Ensure makes sure that the zettelkasten is updated and operational
	Ensure() error
	// GetRoot retrieves the root of the Zettelkasten directory structure
	GetRoot() fs.FS
}

func NewZettel(cfg config.Config) Zettel {
	if cfg.ZettelkastenGitURL != "" {
		return NewGitZettel(cfg)
	} else {
		return NewLocalZettel(cfg)
	}
}
