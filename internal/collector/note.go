package collector

import (
	"log/slog"
	"slices"

	"github.com/luissimas/zettelkasten-exporter/internal/metrics"
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

func CollectNoteMetrics(content []byte) metrics.NoteMetrics {
	links := collectLinks(content)
	linkCount := 0
	for _, v := range links {
		linkCount += v
	}
	return metrics.NoteMetrics{Links: links, LinkCount: linkCount}
}

func collectLinks(content []byte) map[string]int {
	linkKinds := []ast.NodeKind{ast.KindLink, wikilink.Kind}
	reader := text.NewReader(content)
	root := md.Parser().Parse(reader)
	links := make(map[string]int)
	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && slices.Contains(linkKinds, n.Kind()) {
			var target string
			switch v := n.(type) {
			case *ast.Link:
				target = string(v.Destination)
			case *wikilink.Node:
				if v.Embed {
					return ast.WalkContinue, nil
				}
				target = string(v.Target)
			default:
				return ast.WalkContinue, nil
			}

			// TODO: check if target is not a http link
			v, ok := links[target]
			if !ok {
				links[target] = 0
			}
			links[target] = v + 1
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		slog.Error("Error walking note AST", slog.Any("error", err))
	}
	slog.Debug("Collected links", slog.Any("links", links))
	return links
}
