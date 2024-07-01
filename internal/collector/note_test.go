package collector

import (
	"testing"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestCollectNoteMetrics(t *testing.T) {
	data := []struct {
		name     string
		content  string
		expected metrics.NoteMetrics
	}{
		{
			name:    "empty file",
			content: "",
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{},
				LinkCount:     0,
				WordCount:     0,
				BacklinkCount: 0,
			},
		},
		{
			name: "wiki links",
			content: `
[[Link]] some words [[something|another]]

another [[link]]`,
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{"Link": 1, "something": 1, "link": 1},
				LinkCount:     3,
				WordCount:     6,
				BacklinkCount: 0,
			},
		},
		{
			name:    "markdown link",
			content: "[Link](target.md)",
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{"target.md": 1},
				LinkCount:     1,
				WordCount:     1,
				BacklinkCount: 0,
			},
		},
		{
			name:    "repeated links",
			content: "[[target.md|link]] [link](target.md) [[link]]",
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{"target.md": 2, "link": 1},
				LinkCount:     3,
				WordCount:     3,
				BacklinkCount: 0,
			},
		},
		{
			name:    "ignore links to non markdown files",
			content: "![[note.md]] [[test.pdf]] ![[target.png]] ![](another.jpeg) [[link]] [](link)",
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{"link": 2, "note.md": 1},
				LinkCount:     3,
				WordCount:     4,
				BacklinkCount: 0,
			},
		},
		{
			name:    "ignore http links",
			content: "[[one]] [this is an http link](https://go.dev/) [[not/an/http/link]]",
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{"one": 1, "not/an/http/link": 1},
				LinkCount:     2,
				WordCount:     7,
				BacklinkCount: 0,
			},
		},
		{
			name: "mixed links",
			content: `
Ok [Link](target.md).

Another paragraph **bold text** and [[linked]] /test/ [[another|link]].

> Quote in [test](yet-another.md)

A list

- One [[link-unordered.md]]
- Two

Another list:

1. First
2. Second [link](link-ordered.md)`,
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{"target.md": 1, "linked": 1, "another": 1, "yet-another.md": 1, "link-unordered.md": 1, "link-ordered.md": 1},
				LinkCount:     6,
				WordCount:     23,
				BacklinkCount: 0,
			},
		},
		{
			name: "long note",
			content: `
Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.

Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.

Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.

Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.

Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.`,
			expected: metrics.NoteMetrics{
				Links:         map[string]uint{},
				LinkCount:     0,
				WordCount:     525,
				BacklinkCount: 0,
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := CollectNoteMetrics([]byte(d.content))
			assert.Equal(t, d.expected, result)
		})
	}
}
