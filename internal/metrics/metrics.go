package metrics

import "time"

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
	TimeToRead    time.Duration
	BacklinkCount uint
}
