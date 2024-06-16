package collector

import (
	"testing"
	"testing/fstest"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/luissimas/zettelkasten-exporter/internal/storage"
	"github.com/stretchr/testify/assert"
)

func Test_collectMetrics(t *testing.T) {
	fs := fstest.MapFS{
		"zettel/one.md": {Data: []byte(`
---
created-at: "2024-05-29"
---

Testing a note with no links. But there's a [markdown link](./dir1/two.md)

[[./dir1/two.md]]

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
Link to [one](./one.md) and also a full link [[./dir1/dir2/three]] and a [[./dir1/two.md|full link with .md]]
		`)},
		"ignoredir/foo":         {Data: []byte("Foo contents")},
		"ignoredir/bar":         {Data: []byte("Bar contents")},
		"ignoredir/test.md":     {Data: []byte("Test.md contents")},
		"zettel/dir1/ignore.md": {Data: []byte("Ignore.md contents")},
	}
	c := NewCollector(fs, []string{"ignore.md", "ignoredir"}, storage.NewFakeStorage())
	expected := metrics.Metrics{
		NoteCount: 4,
		LinkCount: 8,
		Notes: map[string]metrics.NoteMetrics{
			"zettel/one.md": {
				Links:     map[string]int{"./dir1/two.md": 2},
				LinkCount: 2,
			},
			"zettel/dir1/two.md": {
				Links:     map[string]int{"one": 1},
				LinkCount: 1,
			},
			"zettel/dir1/dir2/three.md": {
				Links:     map[string]int{"one": 1, "two": 1},
				LinkCount: 2,
			},
			"zettel/four.md": {
				Links:     map[string]int{"./one.md": 1, "./dir1/dir2/three": 1, "./dir1/two.md": 1},
				LinkCount: 3,
			},
		},
	}
	metrics, err := c.collectMetrics()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	assert.Equal(t, expected, metrics)
}
