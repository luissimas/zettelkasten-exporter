package zettelkasten

import "io/fs"

// FakeZettelkasten represents a fake implementation of Zettelkasten to be used in tests.
type FakeZettelkasten struct{}

func NewFakeZettelkasten() FakeZettelkasten {
	return FakeZettelkasten{}
}

func (f FakeZettelkasten) Ensure() error {
	return nil
}

func (f FakeZettelkasten) GetRoot() fs.FS {
	return nil
}

func (f FakeZettelkasten) WalkHistory(walkFunc WalkFunc) error {
	return nil
}
