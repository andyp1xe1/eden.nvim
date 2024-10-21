package main

import (
	"bytes"
	"log"

	"garden/appview"
	p "garden/nvim"

	vim "github.com/neovim/go-client/nvim"
	"github.com/yuin/goldmark"
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

	h.ScrollCh() <- int((float64(yCoord)/float64(height))*100) - 25

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
	goldmark.Convert(buf.Bytes(), &htmlBuf)
	return htmlBuf.String()
}
