package collector

import (
	"slices"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/wikilink"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		&wikilink.Extender{},
	),
)

type NoteMetrics struct {
	LinkCount int
}

func CollectNoteMetrics(content []byte) NoteMetrics {
	return NoteMetrics{LinkCount: collectLinkCount(content)}
}

func collectLinkCount(content []byte) int {
	linkKinds := []ast.NodeKind{ast.KindLink, wikilink.Kind}
	reader := text.NewReader(content)
	root := md.Parser().Parse(reader)
	linkCount := 0
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && slices.Contains(linkKinds, n.Kind()) {
			linkCount += 1
		}
		return ast.WalkContinue, nil
	})
	return linkCount
}
