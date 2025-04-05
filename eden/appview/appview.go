package appview

import (
	"fmt"
	"strconv"

	_ "embed"

	web "github.com/webview/webview_go"
)

//go:embed style.css
var style string

const wikiEventListeners = `document.querySelectorAll('a').forEach(l => l.addEventListener('click', e => { e.preventDefault(); window.wikiClick(l.href); console.log(l.href) }))`

type IappView interface {
	Run()
	Destroy()
	LoadDom()
	DocChangedCh() chan<- string
	ScrollCh() chan<- int
	WikiClickCh() <-chan string
}

type AppEvents struct {
	docChangedCh chan string
	scrollCh     chan int
	wikiClickCh  chan string
}

type AppView struct {
	web.WebView
	chans *AppEvents
	title string
}

func makeAppEvents() *AppEvents {
	return &AppEvents{
		docChangedCh: make(chan string),
		scrollCh:     make(chan int),
		wikiClickCh:  make(chan string),
	}
}

func (a AppEvents) DocChangedCh() chan<- string {
	return a.docChangedCh
}

func (a AppEvents) ScrollCh() chan<- int {
	return a.scrollCh
}

func (a AppEvents) WikiClickCh() <-chan string {
	return a.wikiClickCh
}

func MakeAppView(debug bool, title string) *AppView {
	app := &AppView{
		WebView: web.New(debug),
		chans:   makeAppEvents(),
		title:   title,
	}

	return app
}

func (a *AppView) EvenHub() *AppEvents {
	return a.chans
}

func (a *AppView) Run() {

	a.loadDom()

	go func() {
		for {
			select {
			case num := <-a.chans.scrollCh:
				scrollTo(a, num)
			case val := <-a.chans.docChangedCh:
				a.Dispatch(func() {
					updateHTML(a, val)
					a.Eval(wikiEventListeners)
				})
			}
		}
	}()

	a.WebView.Run()
}

func (a *AppView) loadDom() {
	a.SetTitle(a.title)

	a.SetHtml(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<title>Previewer</title>
</head>
<style>%s</style>
<body>
<main>%s</main>
<script>%s</script>
<!--<footer><a href="__prev">previous</a></footer>-->
</body>
</html>`, style, <-a.chans.docChangedCh, wikiEventListeners))

	a.Bind("wikiClick", func(str string) {
		a.chans.wikiClickCh <- str
		a.Eval(`console.log('wiki click consumed')`)
	})
}

func scrollTo(a *AppView, num int) {
	js := fmt.Sprintf(`
dy = (%v/100) * document.body.scrollHeight - (window.innerHeight*0.25);
window.scrollTo({
	top: dy,
	behavior: "smooth"
});`, num)
	a.Dispatch(func() {
		a.Eval(js)
	})
}

func updateHTML(a *AppView, html string) {
	quoted := strconv.Quote(html)
	a.Eval(
		fmt.Sprintf(
			`document.getElementsByTagName('main')[0].innerHTML = %s`,
			quoted))
	a.Eval(wikiEventListeners)
	// a.Eval(`renderMathInElement(document.body)`)
}
