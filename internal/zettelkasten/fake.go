package zettelkasten

import (
	"io/fs"
	"time"
)

// FakeZettelkasten represents a fake implementation of Zettelkasten to be used in tests.
type FakeZettelkasten struct {
	fs fs.FS
}

func NewFakeZettelkasten(fs fs.FS) FakeZettelkasten {
	return FakeZettelkasten{fs: fs}
}

func (f FakeZettelkasten) Ensure() error {
	return nil
}

func (f FakeZettelkasten) GetRoot() fs.FS {
	return f.fs
}

func (f FakeZettelkasten) WalkHistory(walkFunc WalkFunc) error {
	return walkFunc(f.fs, time.Now())
}
