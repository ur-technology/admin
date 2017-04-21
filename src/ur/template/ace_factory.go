package template

import (
	"html/template"
	"ur/file"

	"github.com/yosssi/ace"

	proxy "github.com/yosssi/ace-proxy"
)

type aceFactory struct {
	fs    file.SystemReader
	proxy *proxy.Proxy
}

func NewAceFactory(fs file.SystemReader, baseDir string, dynamicReload bool) Factory {
	return &aceFactory{
		proxy: proxy.New(&ace.Options{
			BaseDir:       baseDir,
			DynamicReload: dynamicReload,
			Asset:         fs.Read,
			DelimLeft:     "<<",
			DelimRight:    ">>",
		}),
		fs: fs,
	}
}

func (a *aceFactory) Template(name string) (*template.Template, error) {
	return a.proxy.Load(name, "", nil)
}
