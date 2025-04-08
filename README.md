# eden.nvim

A makrdown previewer for your digital garden.  
*neovim btw*.

## Installation


### Requirements

At the moment, this is a manual process.
I may make a complementary AUR package in the future.

- install `webkit2gtk` with your package manager
- have go installed
- configure the path so binary is recognized (in `.bashrc` or similar)

```sh
export GOPATH=$HOME/go # you may change this path
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

### lazy

```lua
{
  'andyp1xe1/eden.nvim',
  build = "make",
  config = function()
    require('eden').setup()
  end
}
```

### packer (deprecated?)

```lua
use {
  'andyp1xe1/eden.nvim',
  run = "make",
  config = function()
    require 'eden'.setup()
  end
}
```

### others

~~TODO~~ You're smart, you can figure it out :3

## Usage

`:EdenStart` and `:EdenEnd`

## TODO PROGRESS

- [x] webview with remote updates
- [x] golang neovim client and handlers
- [x] goldmark as parser (and lots of extensions)
- [x] wiki links (goldmak extendsion)
- [x] live scroll
- [x] live html updating
- [x] css styling
- [x] yaml frontmatter support (title and tags)
- [x] two way sync (e.g. navigating wikilinks on the preview opens note in nvim)
- [x] access to offline media (images, etc) (works, but handles paths unsafely)
- [x] mermaid, plantuml and graphviz graphs via kroki servivce + ast transformer
- [ ] latex support (katex? some solution so the rendering is fast while typing)
- [ ] polish css
- [ ] obsidian callouts (via goldmark extension/ast transformer)
- [ ] lua configuration options (custom styling and behavior) + external css support
- [ ] completely factor out / abstract away rendering in a goldmark based module
      (so it is easier to swap out the renderer for a native UI implementation)
- [ ] obsidian like graph mode using [force-graph](https://github.com/vasturiano/force-graph) in case of web rendering, some other way for future native rendering

--- 

## Alternatives

- https://github.com/toppair/peek.nvim

