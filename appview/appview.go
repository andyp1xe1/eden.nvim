package appview

import (
	"fmt"
	"log"
	"strconv"

	web "github.com/webview/webview_go"
)

type AppView interface {
	Run()
	DocChangedCh() chan<- string
	ScrollCh() chan<- int
	LoadDom()
}

type appView struct {
	web.WebView
	title        string
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

	return app
}

func (a *appView) LoadDom() {
	a.SetTitle(a.title)
	a.SetHtml(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<!-- <link rel="stylesheet" href="https://cdn.simplecss.org/simple.min.css">-->
<title>Previewer</title>
</head>
<style>%s</style>
<body>%s</body>
</html>`, style, <-a.docChangedCh))
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
			js := fmt.Sprintf(`requestAnimationFrame(function () {
	percentage = (%v/100) * document.body.scrollHeight - 50
	window.scrollTo({
		top: percentage,
		behavior: "smooth"
	})})`, num)
			log.Println(js)
			a.Dispatch(func() {
				a.Eval(js)
			})
		}
	}()
	go func() {
		for val := range a.docChangedCh {
			a.Dispatch(func() {
				quoted := strconv.Quote(val)
				a.Eval(
					fmt.Sprintf(`requestAnimationFrame(function () {
document.getElementsByTagName('body')[0].innerHTML = %s
})`, quoted))
			})
		}
	}()
	a.WebView.Run()
}

const style = `
/* Global variables. */
:root,
::backdrop {
  /* Set sans-serif & mono fonts */
  --sans-font: -apple-system, BlinkMacSystemFont, "Avenir Next", Avenir,
    "Nimbus Sans L", Roboto, "Noto Sans", "Segoe UI", Arial, Helvetica,
    "Helvetica Neue", sans-serif;
  --mono-font: "Fantasque Sans Mono", Consolas, Menlo, Monaco, "Andale Mono", "Ubuntu Mono", monospace;
  --standard-border-radius: 5px;

  /* Gruvbox Light theme */
  --bg: #fbf1c7; /* light background */
  --accent-bg: #ebdbb2;
  --text: #3c3836; /* dark text */
  --text-light: #7c6f64; /* soft dark text */
  --border: #928374; /* light gray border */
  --accent: #d65d0e; /* orange accent */
  --accent-hover: #fe8019; /* brighter orange */
  --accent-text: var(--bg);
  --code: #b16286; /* purple/pink for code */
	--code-bg: #ebdbb2; /* code background */
  --preformatted: #458588; /* aqua for preformatted text */
  --marked: #fabd2f; /* yellow for highlights */
  --disabled: #f2e5bc; /* light disabled state */
}

/* Dark theme */
@media (prefers-color-scheme: dark) {
  :root,
  ::backdrop {
    color-scheme: dark;
    --bg: #282828; /* dark background */
    --accent-bg: #3c3836;
    --text: #ebdbb2; /* light text */
    --text-light: #a89984; /* soft light text */
    --accent: #d65d0e; /* orange accent */
    --accent-hover: #fe8019; /* brighter orange */
    --accent-text: var(--bg);
    --code: #d3869b; /* purple/pink for code */
		--code-bg: #3c3836; /* code background */
    --preformatted: #83a598; /* aqua for preformatted text */
    --disabled: #504945; /* dark disabled state */
  }
  /* Add a bit of transparency so light media isn't so glaring in dark mode */
  img,
  video {
    opacity: 0.8;
  }
}

/* Reset box-sizing */
*, *::before, *::after {
  box-sizing: border-box;
}

html, body {
  height: 100%;
}

html {
  /* Set the font globally */
  font-family: var(--sans-font);
  scroll-behavior: smooth;
}

/* Make the body a nice central block */
body {
  color: var(--text);
  background-color: var(--bg);
  font-size: 1.15rem;
  line-height: 1.5;
  display: grid;
  grid-template-columns: 1fr min(45rem, 90%) 1fr;
  grid-template-rows: min-content 1fr auto;
  margin: 0;
}
body > * {
  grid-column: 2;
}

main {
  padding-top: 1.5rem;
}

body > footer {
  margin-top: 4rem;
  padding: 2rem 1rem;
  color: var(--text-light);
  font-size: 0.9rem;
  text-align: center;
  border-top: 1px solid var(--border);
}

/* Format headers */
h1 {
  font-size: 3rem;
}

h2 {
  font-size: 2.6rem;
  margin-top: 3rem;
}

h3 {
  font-size: 2rem;
  margin-top: 3rem;
}

h4 {
  font-size: 1.44rem;
}

h5 {
  font-size: 1.15rem;
}

h6 {
  font-size: 0.96rem;
}

p {
	margin: 0;
  padding: 1.5rem 0;
}

/* Prevent long strings from overflowing container */
p, h1, h2, h3, h4, h5, h6 {
  overflow-wrap: break-word;
}

/* Fix line height when title wraps */
h1,
h2,
h3 {
  line-height: 1.1;
}

/* Reduce header size on mobile */
@media only screen and (max-width: 720px) {
  h1 {
    font-size: 2.5rem;
  }

  h2 {
    font-size: 2.1rem;
  }

  h3 {
    font-size: 1.75rem;
  }

  h4 {
    font-size: 1.25rem;
  }
}

/* Format links*/
a,
a:visited {
  color: var(--accent);
	pointer-events: none;
}

a:hover {
  text-decoration: none;
}

/* Format tables */
table {
  border-collapse: collapse;
  margin: 1.5rem 0;
}

figure > table {
  width: max-content;
  margin: 0;
}

td,
th {
  border: 1px solid var(--border);
  text-align: start;
  padding: 0.5rem;
}

th {
  background-color: var(--accent-bg);
  font-weight: bold;
}

tr:nth-child(even) {
  /* Set every other cell slightly darker. Improves readability. */
  background-color: var(--accent-bg);
}

table caption {
  font-weight: bold;
  margin-bottom: 0.5rem;
}

/* Misc body elements */
hr {
  border: none;
  height: 1px;
  background: var(--border);
  margin: 1rem auto;
}

mark {
  padding: 2px 5px;
  border-radius: var(--standard-border-radius);
  background-color: var(--marked);
  color: black;
}

mark a {
  color: #0d47a1;
}

img,
video {
  max-width: 100%;
  height: auto;
  border-radius: var(--standard-border-radius);
}

figure {
  margin: 0;
  display: block;
  overflow-x: auto;
}

figure > img,
figure > picture > img {
  display: block;
  margin-inline: auto;
}

figcaption {
  text-align: center;
  font-size: 0.9rem;
  color: var(--text-light);
  margin-block: 1rem;
}

blockquote {
	margin: 0;
	background: var(--accent-bg);
	border: 1px solid var(--border);

  border-left: 0.5rem solid var(--accent);
  padding: 0.3rem 1rem;
  
  /*border-radius: var(--standard-border-radius);*/
  color: var(--text-light);
  /*font-style: italic;*/
}

cite {
  font-size: 0.9rem;
  color: var(--text-light);
  font-style: normal;
}

dt {
    color: var(--text-light);
}

/* Use mono font for code elements */
code,
pre,
pre span,
kbd,
samp {
  font-family: var(--mono-font);
  color: var(--code);
}

kbd {
  color: var(--preformatted);
  border: 1px solid var(--preformatted);
  border-bottom: 3px solid var(--preformatted);
  border-radius: var(--standard-border-radius);
  padding: 0.1rem 0.4rem;
}


pre {
  border: 1px solid var(--border);
  background-color: var(--code-bg);
  white-space: pre;
  tab-size: 2;
  margin: 1rem 0;
  padding: 1rem;
  font-size: 1rem;
  line-height: 1.4;
  color: var(--preformatted);
  border-radius: var(--standard-border-radius);
}

pre > code {
  color: inherit;
  background: none;
  border: none;
  padding: 0;
  line-height: inherit;
  display: block;
}

/* Superscript & Subscript */
/* Prevent scripts from affecting line-height. */
sup, sub {
  vertical-align: baseline;
  position: relative;
}

sup {
  top: -0.4em;
}

sub { 
  top: 0.3em; 
}
`
