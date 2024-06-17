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
				Links:     map[string]int{},
				LinkCount: 0,
			},
		},
		{
			name:    "wiki links",
			content: "[[Link]]aksdjf[[something|another]]\n[[link]]",
			expected: metrics.NoteMetrics{
				Links:     map[string]int{"Link": 1, "something": 1, "link": 1},
				LinkCount: 3,
			},
		},
		{
			name:    "markdown link",
			content: "[Link](target.md)",
			expected: metrics.NoteMetrics{
				Links:     map[string]int{"target.md": 1},
				LinkCount: 1,
			},
		},
		{
			name:    "mixed links",
			content: "okok[Link](target.md)\n**ddk**[[linked]]`test`[[another|link]]\n\n[test](yet-another.md)",
			expected: metrics.NoteMetrics{
				Links:     map[string]int{"target.md": 1, "linked": 1, "another": 1, "yet-another.md": 1},
				LinkCount: 4,
			},
		},
		{
			name:    "repeated links",
			content: "[[target.md|link]]\n[link](target.md)\n[[link]]",
			expected: metrics.NoteMetrics{
				Links:     map[string]int{"target.md": 2, "link": 1},
				LinkCount: 3,
			},
		},
		{
			name:    "ignore embeddedlinks",
			content: "![[target.png]]\n!()[another.jpeg]\n[[link]]",
			expected: metrics.NoteMetrics{
				Links:     map[string]int{"link": 1},
				LinkCount: 1,
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := CollectNoteMetrics([]byte(d.content))
			assert.Equal(t, d.expected.Links, result.Links)
		})
	}
}
