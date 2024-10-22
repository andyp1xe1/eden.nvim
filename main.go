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
	if err := md.Convert(buf.Bytes(), &htmlBuf, parser.WithContext(ctx)); err != nil {
		return fmt.Sprintf(`<h1>Error</h1>
<pre>%v</pre>
`, err.Error())
	}
	metadata := meta.Get(ctx)
	if title, ok := metadata["title"]; !ok {
		return htmlBuf.String()
	} else {
		return fmt.Sprintf(`<h1 styles="margin: 0 auto; width: fit-content;" >%s<h1>%s`, title, &htmlBuf)
	}
}
