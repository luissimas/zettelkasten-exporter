package collector

import (
	"log/slog"
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
	Links map[string]int
}

func CollectNoteMetrics(content []byte) NoteMetrics {
	return NoteMetrics{Links: collectLinks(content)}
}

func collectLinks(content []byte) map[string]int {
	linkKinds := []ast.NodeKind{ast.KindLink, wikilink.Kind}
	reader := text.NewReader(content)
	root := md.Parser().Parse(reader)
	links := make(map[string]int)
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && slices.Contains(linkKinds, n.Kind()) {
			var target string
			switch v := n.(type) {
			case *ast.Link:
				target = string(v.Destination)
			case *wikilink.Node:
				target = string(v.Target)
			}

			v, ok := links[target]
			if !ok {
				links[target] = 0
			}
			links[target] = v + 1
		}
		return ast.WalkContinue, nil
	})
	slog.Info("Collected links", slog.Any("links", links))
	return links
}
