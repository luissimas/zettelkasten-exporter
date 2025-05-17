package metrics

// ZettelkastenMetrics represents the aggregated metrics of a Zettelkasten.
type ZettelkastenMetrics struct {
	NoteCount uint
	LinkCount uint
	WordCount uint
	Notes     map[string]NoteMetrics
}

// NoteMetrics represents the metrics of a single Zettelkasten note.
type NoteMetrics struct {
	Links         map[string]uint
	LinkCount     uint
	WordCount     uint
	BacklinkCount uint
}
