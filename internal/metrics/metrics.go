package metrics

type Metrics struct {
	NoteCount uint
	LinkCount uint
	WordCount uint
	Notes     map[string]NoteMetrics
}

type NoteMetrics struct {
	Links         map[string]uint
	LinkCount     uint
	WordCount     uint
	BacklinkCount uint
}
