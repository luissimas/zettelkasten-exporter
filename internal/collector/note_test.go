package collector

import "testing"

func TestCollectNoteMetrics(t *testing.T) {
	data := []struct {
		name     string
		content  string
		expected NoteMetrics
	}{
		{"empty file", "", NoteMetrics{LinkCount: 0}},
		{"file with a single link", "[[Link]]", NoteMetrics{LinkCount: 1}},
		{"file with multiple links", "[[Link]]aksdjf[[anotherlink]]\n[[link]]", NoteMetrics{LinkCount: 3}},
		{"wikilink dividers", "[[something|another]]\n\n[[link]]\n[[382dlk djfs link|yeah]]", NoteMetrics{LinkCount: 3}},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := CollectNoteMetrics([]byte(d.content))
			if result != d.expected {
				t.Errorf("Expected %v, got %v", d.expected, result)
			}
		})
	}
}
