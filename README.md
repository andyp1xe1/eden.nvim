# eden.nvim

A makrdown previewer for your digital garden.  
*neovim btw*.

## Installation

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
- [ ] access to offline media (images, etc)
- [ ] obsidian callouts
- [ ] latex support
- [ ] lua configuration options (custom styling and behavior)
- [ ] markdown formatter (something that uses goldmark)

--- 

## Alternatives

- https://github.com/toppair/peek.nvim
