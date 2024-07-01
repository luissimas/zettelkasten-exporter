package collector

import (
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
	"unicode"

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
	noteMetrics := metrics.NoteMetrics{
		Links:         make(map[string]uint),
		LinkCount:     0,
		WordCount:     0,
		BacklinkCount: 0,
	}
	reader := text.NewReader(content)
	root := md.Parser().Parse(reader)
	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		linkTarget := ""

		switch v := n.(type) {
		case *ast.Link:
			linkTarget = string(v.Destination)
		case *wikilink.Node:
			linkTarget = string(v.Target)
		case *ast.Paragraph, *ast.ListItem:
			text := string(n.Text(content))
			fields := strings.FieldsFunc(string(text), func(r rune) bool { return unicode.IsSpace(r) || r == '\n' })
			noteMetrics.WordCount += uint(len(fields))
		default:
			return ast.WalkContinue, nil
		}

		if !isNoteTarget(linkTarget) {
			return ast.WalkContinue, nil
		}

		v, ok := noteMetrics.Links[linkTarget]
		if !ok {
			noteMetrics.Links[linkTarget] = 0
		}
		noteMetrics.Links[linkTarget] = v + 1
		return ast.WalkContinue, nil
	})
	if err != nil {
		slog.Error("Error walking note AST", slog.Any("error", err))
	}
	for _, linkCount := range noteMetrics.Links {
		noteMetrics.LinkCount += linkCount
	}
	return noteMetrics
}

// isNoteTarget determines whether a link target points to a markdown note.
func isNoteTarget(target string) bool {
	// Empty strings are not valid targets
	if target == "" {
		return false
	}

	// Check if target is a URL
	u, err := url.Parse(target)
	isUrl := err == nil && u.Scheme != "" && u.Host != ""
	if isUrl {
		return false
	}

	// Check if target is either a markdown file or has no extension
	extension := filepath.Ext(target)
	isNoteTarget := extension == "" || extension == ".md"
	return isNoteTarget
}
