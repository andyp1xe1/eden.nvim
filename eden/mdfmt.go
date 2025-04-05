package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/wikilink"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// TODO:
// - [ ] insdead of appending html make an extension for goldmark?

type plainResolver struct{}

func (plainResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	u := n.Target
	if len(n.Fragment) > 0 {
		u = append(u, '#')
		u = append(u, n.Fragment...)
	}
	if n.Embed {
		return append([]byte{'!'}, u...), nil
	}
	return u, nil
}

func parseLines(lines [][]byte) string {
	var buf bytes.Buffer
	var htmlBuf bytes.Buffer

	for i, line := range lines {
		buf.Write(line)
		if i < len(lines)-1 {
			buf.WriteByte('\n')
		}
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			extension.GFM,
			extension.Table,
			extension.Typographer,
			extension.Strikethrough,
			extension.Footnote,
			&wikilink.Extender{Resolver: plainResolver{}},
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	ctx := parser.NewContext()

	err := md.Convert(
		buf.Bytes(),
		&htmlBuf,
		parser.WithContext(ctx))

	if err != nil {
		return fmt.Sprintf(`<h1>Error</h1>
<pre>%v</pre>`, err.Error())
	}

	var front string
	metadata := meta.Get(ctx)
	if title, ok := metadata["title"]; ok {
		front += fmtTitle(title)
	}
	if tags, ok := metadata["tags"]; ok {
		front += fmtTags(tags)
	}
	return front + "\n" + prevHeader + "\n" + htmlBuf.String()
}

func fmtTitle(title interface{}) string {
	var html string
	if str, ok := title.(string); ok {
		html = fmt.Sprintf(
			`<h1 class="title">%s</h1>`, str)
	}
	return html
}

func fmtTags(tags interface{}) string {
	var html string
	if list, ok := tags.([]interface{}); ok {
		for _, li := range list {
			html += fmt.Sprintf(
				`<li><a href="__tag#%s">#%s</a></li>`, li, li)
		}
		return fmt.Sprintf(
			`<ul class="tags" >%s</ul>`, html)
	}
	return html
}
