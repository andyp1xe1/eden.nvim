package main

import (
	web "github.com/webview/webview_go"
)

type AppView struct {
 web.WebView
 Title, Path string
}

func NewAppView(debug bool, title, path string) *AppView{
	webView := web.New(debug)
	return &AppView{webView, title, path}
}

func (a *AppView) Init() {
	a.SetTitle(a.Title)
	a.Navigate(a.Path)
  a.SetHtml(`
		<!DOCTYPE html>
		<html>
		<body>
			<script>
				window.scrollToTop = function() {
					window.scrollTo(0, 0);
				}

				window.scrollToBottom = function() {
					window.scrollTo(0, document.body.scrollHeight);
				}
			</script>
		</body>
		</html>
	`)
}

