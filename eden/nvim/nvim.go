package nvim

import (
	"log"
	"os"

	vim "github.com/neovim/go-client/nvim"
)

type HandlerMap map[string]Handler
type Handler interface{} //func(v *vim.Nvim, args []string)

type Conf struct {
	Name     string
	Handlers HandlerMap
}

type NvimPlugin interface {
	Serve() error
	Vim() *vim.Nvim
}

type nvimPlugin struct {
	*vim.Nvim
	conf Conf
}

func (p *nvimPlugin) Vim() *vim.Nvim {
	return p.Nvim
}

func Setup(c Conf) (NvimPlugin, error) {
	v, err := newStdio(log.Printf)
	if err != nil {
		return nil, err
	}
	plugin := &nvimPlugin{v, c}
	plugin.registerHandlers()
	return plugin, nil
}

func (np *nvimPlugin) Serve() error {
	return np.Nvim.Serve()
}

func (p *nvimPlugin) registerHandlers() {
	for name, fun := range p.conf.Handlers {
		p.RegisterHandler(name, fun)
	}
}

func newStdio(logf func(string, ...interface{})) (*vim.Nvim, error) {
	log.SetFlags(0)
	stdout := os.Stdout
	os.Stdout = os.Stderr

	v, err := vim.New(os.Stdin, stdout, stdout, logf)
	if err != nil {
		return nil, err
	}
	return v, nil
}
