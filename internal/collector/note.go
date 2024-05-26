package collector

import (
	"regexp"
)

type NoteMetrics struct {
	LinkCount int
}

func CollectNoteMetrics(content []byte) NoteMetrics {
	return NoteMetrics{LinkCount: collectLinkCount(content)}
}

func collectLinkCount(content []byte) int {
	r, _ := regexp.Compile(`\[\[[^]]+\]\]`)
	matches := r.FindAll(content, -1)
	return len(matches)
}
