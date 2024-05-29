package collector

import "testing"

func TestCollectNoteMetrics(t *testing.T) {
	data := []struct {
		name     string
		content  string
		expected NoteMetrics
	}{
		{"empty file", "", NoteMetrics{LinkCount: 0}},
		{"file with a single wiki link", "[[Link]]", NoteMetrics{LinkCount: 1}},
		{"file with multiple wiki links", "[[Link]]aksdjf[[anotherlink]]\n[[link]]", NoteMetrics{LinkCount: 3}},
		{"wikilink dividers", "[[something|another]]\n\n[[link]]\n[[382dlk djfs link|yeah]]", NoteMetrics{LinkCount: 3}},
		{"file with markdown link", "[Link](target.md)", NoteMetrics{LinkCount: 1}},
		{"file with multiple links", "okok[Link](target.md)\n**ddk**[[linked]]`test`[[anothe|link]]\n\n[test](yet-another.md)", NoteMetrics{LinkCount: 4}},
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
