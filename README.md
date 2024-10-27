# eden.nvim

A makrdown previewer for your digital garden. *In neovim ofcourse.*

## Installation

### packer

```lua
use {
  'andyp1xe1/eden.nvim',
  run = "make",
  config = function()
    require 'eden'.setup()
  end
}
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
- [ ] access to offline media (images, etc)
- [ ] obsidian callouts
- [ ] latex support
- [ ] lua configuration options (custom styling and behavior)
- [ ] two way sync (e.g. navigating links on the preview opens note in nvim)
- [ ] markdown formatter (something that uses goldmark)

--- 

## Alternatives

- https://github.com/toppair/peek.nvim
