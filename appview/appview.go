package appview

import (
	"fmt"
	"log"

	web "github.com/webview/webview_go"
)

type AppView interface {
	Run()
	DocChangedCh() chan<- string
	ScrollCh() chan<- int
}

type appView struct {
	web.WebView
	title string

	docChangedCh chan string
	scrollCh     chan int
}

func Setup(debug bool, title string) *appView {
	app := &appView{
		WebView:      web.New(debug),
		title:        title,
		docChangedCh: make(chan string),
		scrollCh:     make(chan int),
	}

	app.SetTitle(title)
	app.SetHtml(`
		<!DOCTYPE html>
		<html>
		<body>
			<h1>Loading</h1>
			<script>
			<!-- Potential Functions Defined -->
			</script>
		</body>
		</html>
	`)
	return app
}

func (a *appView) DocChangedCh() chan<- string {
	return a.docChangedCh
}

func (a *appView) ScrollCh() chan<- int {
	return a.scrollCh
}

func (a *appView) Run() {
	go func() {
		for num := range a.scrollCh {
			js := fmt.Sprintf(`
	percentage = (%v/100) * document.body.scrollHeight
	window.scrollTo({
		top: percentage,
		behavior: "smooth"
	})`, num)
			log.Println(js)
			a.Eval(js)
		}
	}()
	go func() {
		for val := range a.docChangedCh {
			a.SetHtml(val)
		}
	}()
	a.WebView.Run()
}
