package exporter

import (
	"context"
	"testing"
	"testing/fstest"
	"time"

	"github.com/luissimas/zettelkasten-exporter/internal/config"
	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
	"github.com/luissimas/zettelkasten-exporter/internal/zettelkasten"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	fs := fstest.MapFS{
		"zettel/one.md": {Data: []byte(`
---
created-at: "2024-05-29"
---

Testing a note with no links. But there's a [markdown link](./dir1/two.md)

[[two]]

![[./image.png]]
		`)},
		"zettel/dir1/two.md": {Data: []byte(`
---
created-at: "2024-05-29"
---

This note links to [[one]]

![](./image.png)
		`)},
		"zettel/dir1/dir2/three.md": {Data: []byte(`
---
created-at: "2024-05-29"
---

Links to [[one]] but also to [[two|two with an alias]]
		`)},
		"zettel/four.md": {Data: []byte(`
---
created-at: "2024-05-29"
---
Link to [one](one.md) and also a full link [[./dir1/dir2/three]] and a [[dir1/two.md|full link with .md]]

- item
-
- another item
		`)},
		"ignoredir/foo":         {Data: []byte("Foo contents")},
		"ignoredir/bar":         {Data: []byte("Bar contents")},
		"ignoredir/test.md":     {Data: []byte("Test.md contents")},
		"zettel/dir1/ignore.md": {Data: []byte("Ignore.md contents")},
	}
	fakeStorage := storage.NewFakeStorage()
	exporter := NewExporter(config.Config{IgnoreFiles: []string{"ignore.md", "ignoredir"}, CollectionInterval: time.Millisecond * 10}, zettelkasten.NewFakeZettelkasten(fs), &fakeStorage)
	expected := metrics.ZettelkastenMetrics{
		NoteCount: 4,
		LinkCount: 8,
		WordCount: 47,
		Notes: map[string]metrics.NoteMetrics{
			"one": {
				Links:         map[string]uint{"two": 2},
				LinkCount:     2,
				WordCount:     13,
				BacklinkCount: 3,
			},
			"two": {
				Links:         map[string]uint{"one": 1},
				LinkCount:     1,
				WordCount:     6,
				BacklinkCount: 4,
			},
			"three": {
				Links:         map[string]uint{"one": 1, "two": 1},
				LinkCount:     2,
				WordCount:     10,
				BacklinkCount: 1,
			},
			"four": {
				Links:         map[string]uint{"one": 1, "three": 1, "two": 1},
				LinkCount:     3,
				WordCount:     18,
				BacklinkCount: 0,
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	err := exporter.Start(ctx)

	require.NoError(t, err)
	assert.Equal(t, len(fakeStorage.Metrics), 10)
	for _, metric := range fakeStorage.Metrics {
		assert.Equal(t, expected, metric)
	}
}
