package template

import "html/template"

type Factory interface {
	Template(name string) (*template.Template, error)
}
