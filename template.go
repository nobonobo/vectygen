package main

import "text/template"

var templ = template.Must(template.New("").Parse(`package {{.PkgName}}

import (
{{range $v, $_ := .StdImports}}{{printf "\t%q\n" $v}}{{end -}}
{{if .StdImports}}{{printf "\n"}}{{end -}}
{{range $v, $_ := .Imports}}{{printf "\t%q\n" $v}}{{end -}}
)

// New{{.ComponentName}} ...
func New{{.ComponentName}}(d map[string]func(*vecty.Event)) *{{.ComponentName}} {
	return &{{.ComponentName}}{
		dispatcher: d,
	}
}

// {{.ComponentName}} ...
type {{.ComponentName}} struct{
	vecty.Core
	dispatcher map[string]func(*vecty.Event)
}

// Render ...
func (c *{{.ComponentName}}) Render() vecty.ComponentOrHTML {
	return {{.Generated}}
}

{{range $event, $method := .Methods -}}
// {{$method}} ...
func (c *{{$.ComponentName}}) {{$method}}(event *vecty.Event) {
	f, ok := c.dispatcher["{{$method}}"]
	if !ok {
		panic("unknown func: \"{{$method}}\"")
	}
	f(event)
}
{{end}}
`))
