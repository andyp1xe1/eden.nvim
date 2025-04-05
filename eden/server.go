package main

import (
	"fmt"
	"log"
	"net/url"

	nvim "github.com/andyp1xe1/eden.nvim/eden/nvim"
	vim "github.com/neovim/go-client/nvim"
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
