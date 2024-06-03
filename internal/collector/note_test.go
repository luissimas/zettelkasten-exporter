package collector

import (
	"maps"
	"testing"
)

func TestCollectNoteMetrics(t *testing.T) {
	data := []struct {
		name     string
		content  string
		expected NoteMetrics
	}{
		{
			name:     "empty file",
			content:  "",
			expected: NoteMetrics{Links: map[string]int{}},
		},
		{
			name:     "wiki links",
			content:  "[[Link]]aksdjf[[something|another]]\n[[link]]",
			expected: NoteMetrics{Links: map[string]int{"Link": 1, "something": 1, "link": 1}},
		},
		{
			name:     "markdown link",
			content:  "[Link](target.md)",
			expected: NoteMetrics{Links: map[string]int{"target.md": 1}},
		},
		{
			name:     "mixed links",
			content:  "okok[Link](target.md)\n**ddk**[[linked]]`test`[[another|link]]\n\n[test](yet-another.md)",
			expected: NoteMetrics{Links: map[string]int{"target.md": 1, "linked": 1, "another": 1, "yet-another.md": 1}},
		},
		{
			name:     "repeated links",
			content:  "[[target.md|link]]\n[link](target.md)\n[[link]]",
			expected: NoteMetrics{Links: map[string]int{"target.md": 2, "link": 1}},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := CollectNoteMetrics([]byte(d.content))
			if !maps.Equal(result.Links, d.expected.Links) {
				t.Errorf("Expected %v, got %v", d.expected, result)
			}
		})
	}
}
