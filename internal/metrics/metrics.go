package metrics

type Metrics struct {
	NoteCount int
	LinkCount int
	Notes     map[string]NoteMetrics
}

type NoteMetrics struct {
	Links     map[string]int
	LinkCount int
}
