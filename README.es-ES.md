

# eden.nvim

Un visor de Markdown para tu jardín digital.  
*neovim, por cierto*.

## Instalación


### Requisitos

Por el momento, este es un proceso manual.
Puede que cree un paquete complementario para AUR en el futuro.

- instala `webkit2gtk` con tu gestor de paquetes
- tener Go instalado
- configurar la ruta para que el binario sea reconocido (en `.bashrc` o similar)

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

### packer (¿deprecado?)

```lua
use {
  'andyp1xe1/eden.nvim',
  run = "make",
  config = function()
    require 'eden'.setup()
  end
}
```

### otros

~~TODO~~ Eres inteligente, puedes averiguarlo tú mismo :3

## Uso

`:EdenStart` y `:EdenEnd`

## PROGRESO DEL TODO

- [x] webview con actualizaciones remotas
- [x] cliente de Neovim en Go y manejadores
- [x] goldmark como analizador (y muchas extensiones)
- [x] enlaces wiki (extensión de goldmark)
- [x] desplazamiento en vivo
- [x] actualización de HTML en vivo
- [x] estilos CSS
- [x] soporte para frontmatter YAML (título y etiquetas)
- [x] sincronización bidireccional (p. ej., navegar por enlaces wiki en la vista previa abre la nota en nvim)
- [x] acceso a medios sin conexión (imágenes, etc.) (funciona, pero maneja las rutas de forma insegura)
- [x] gráficos de mermaid, plantuml y graphviz a través del servicio kroki + transformador AST
- [ ] soporte para LaTeX (¿katex? alguna solución para que el renderizado sea rápido mientras se escribe)
- [ ] pulir el CSS
- [ ] callouts de Obsidian (vía extensión de goldmark/transformador AST)
- [ ] opciones de configuración en Lua (estilos y comportamiento personalizados) + soporte para CSS externo
- [ ] extraer por completo / abstraer el renderizado en un módulo basado en goldmark
      (para que sea más fácil cambiar el renderizador por una implementación de interfaz de usuario nativa)
- [ ] modo de grafo al estilo de Obsidian usando [force-graph](https://github.com/vasturiano/force-graph) en caso de renderizado web, y alguna otra forma para el renderizado nativo futuro

--- 

## Alternativas

- https://github.com/toppair/peek.nvim
