package main

import (
	"bytes"
	"fmt"
	"log"

	"garden/appview"
	p "garden/nvim"

	vim "github.com/neovim/go-client/nvim"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/wikilink"
)

type Handler struct {
	appview.AppView
}

func main() {
	log.SetFlags(0)

	app := appview.Setup(
		true,
		"Markdown Preview",
	)

	handler := Handler{app}

	go handler.serve()
	app.LoadDom()

	app.Run()
	app.Destroy()
}

func (h Handler) serve() {
	plugin, err := p.Setup(p.Conf{
		Name: "Markdown Preview",
		Handlers: p.HandlerMap{
			"text_changed": h.onTextChanged,
			"scroll":       h.onScroll,
			"enter":        h.onBufEnter,
		},
	})
	if err != nil {
		panic(err)
	}
	if err := plugin.Serve(); err != nil {
		panic(err)
	}
}

func (h Handler) onScroll(v *vim.Nvim) error {
	vec, err := v.WindowCursor(0)
	if err != nil {
		return err
	}

	height, err := v.BufferLineCount(0)
	if err != nil {
		return err
	}

	yCoord := vec[0]

	h.ScrollCh() <- int((float64(yCoord) / float64(height)) * 100)

	return nil
}

// This thing is here just in case I may need it in the future
func (h Handler) onBufEnter(v *vim.Nvim) error { return nil }

func (h Handler) onTextChanged(v *vim.Nvim) error {
	buf, err := v.CurrentBuffer()
	if err != nil {
		return err
	}
	lines, err := v.BufferLines(buf, 0, -1, true)
	if err != nil {
		return err
	}
	h.DocChangedCh() <- parseLines(lines)
	return nil
}

// TODO:
// - move in another file or package?
// - insdead of appending html make an extension for goldmark?
// - resolve the links to something meaingfull
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
			&wikilink.Extender{},
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
	return front + "\n" + htmlBuf.String()
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
				`<li><a href="#">#%s</a></li>`, li)
		}
		return fmt.Sprintf(
			`<ul class="tags" >%s</ul>`, html)
	}
	return html
}
