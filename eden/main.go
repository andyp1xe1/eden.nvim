package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"

	"github.com/andyp1xe1/eden.nvim/eden/appview"
	nvim "github.com/andyp1xe1/eden.nvim/eden/nvim"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	vim "github.com/neovim/go-client/nvim"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/wikilink"
)

type EventHub interface {
	DocChangedCh() chan<- string
	ScrollCh() chan<- int
	WikiClickCh() <-chan string
}

type Plugin interface {
	Serve() error
	Vim() *vim.Nvim
}

type PluginServer struct {
	hub    EventHub
	plugin Plugin
}

func main() {
	log.SetFlags(0)

	app := appview.MakeAppView(
		true,
		"Markdown Preview",
	)

	server := MakeServer(app.EvenHub())
	server.Serve()

	app.Run()
	app.Destroy()
}

func Handler(hub EventHub, fn func(h EventHub, v *vim.Nvim) error) func(v *vim.Nvim) error {
	return func(v *vim.Nvim) error {
		return fn(hub, v)
	}
}

func MakeServer(hub EventHub) *PluginServer {
	return &PluginServer{hub: hub}
}

func (s PluginServer) Serve() {
	var err error
	if s.plugin, err = nvim.Setup(nvim.Conf{
		Name: "Markdown Preview",
		Handlers: nvim.HandlerMap{
			"text_changed": Handler(s.hub, onTextChanged),
			"scroll":       Handler(s.hub, onScroll),
			"enter":        Handler(s.hub, onBufEnter),
		},
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := s.plugin.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	go s.TalkBack()

}

func (s PluginServer) TalkBack() {
	for target := range s.hub.WikiClickCh() {
		clickWiki(s.plugin.Vim(), target)
	}
	// TODO: use `for-select` for more future takk backs
}

func clickWiki(v *vim.Nvim, target string) {
	target, err := url.PathUnescape(target)
	if err != nil {
		v.WritelnErr("Error: " + err.Error())
	}
	v.Command("let @/ = ''")
	v.Command(fmt.Sprintf(`call search('\[\[%s')`, target))
	v.Command(`normal gd`)
	// v.Command(`nohlsearch|redraw`)
}

func onScroll(h EventHub, v *vim.Nvim) error {
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
func onBufEnter(h EventHub, v *vim.Nvim) error { return nil }

func onTextChanged(h EventHub, v *vim.Nvim) error {
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
// - [ ] move in another file or package?
// - [ ] insdead of appending html make an extension for goldmark?
// - [x] resolve the ~links~ wikilinks to something meaingfull
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
	return front + "\n" + htmlBuf.String()
}

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
