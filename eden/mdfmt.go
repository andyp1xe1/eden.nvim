package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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
			parser.WithASTTransformers(
				util.PrioritizedValue{
					Value:    &KrokiTransformer{},
					Priority: 200,
				},
				util.PrioritizedValue{
					Value:    &LocalImageTransformer{Port: "6969"},
					Priority: 100,
				},
			),
		))

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

var diagramTypes = map[string]string{
	"dot":      "graphviz",
	"graphviz": "graphviz",
	"mermaid":  "mermaid",
	"plantuml": "plantuml",
	"puml":     "plantuml",
}

func allDiagramTypesRe() string {
	var types []string
	for k := range diagramTypes {
		types = append(types, k)
	}
	return strings.Join(types, "|")
}

type KrokiTransformer struct{}

func (KrokiTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindFencedCodeBlock {
			return ast.WalkContinue, nil
		}

		block := n.(*ast.FencedCodeBlock)
		lang := string(block.Language(reader.Source()))

		dtype, ok := diagramTypes[lang]
		if !ok {
			return ast.WalkContinue, nil
		}

		content := block.Lines().Value(reader.Source())
		url, err := makeURL([]byte(dtype), []byte("svg"), content)
		if err != nil {
			return ast.WalkContinue, err
		}

		link := ast.NewLink()
		link.Destination = url
		link.Title = []byte(dtype + " diagram")
		img := ast.NewImage(link)

		block.Parent().ReplaceChild(block.Parent(), block, img)

		return ast.WalkSkipChildren, nil
	})
}

type LocalImageTransformer struct {
	Port string
}

func (t *LocalImageTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	// Visit all nodes in the AST
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}

		image := n.(*ast.Image)
		destURL := string(image.Destination)

		if !strings.HasPrefix(destURL, "/") && !strings.HasPrefix(destURL, "./") {
			return ast.WalkContinue, nil
		}

		if strings.HasPrefix(destURL, "./") {
			destURL = destURL[1:]
		}

		newURL := fmt.Sprintf("http://localhost:%s%s", t.Port, destURL)
		image.Destination = []byte(newURL)

		return ast.WalkContinue, nil
	})
}

// func addPortToImages(html string, port string) string {
// 	re := regexp.MustCompile(`<img src="\.?(/[^"]*)"`)
// 	return re.ReplaceAllString(html, `<img src="http://localhost:`+port+`$1"`)
// }

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
