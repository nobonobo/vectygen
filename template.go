package main

import "text/template"

var templ = template.Must(template.New("").Parse(`package {{.PkgName}}

import (
{{range $v := .StdImports}}{{printf "\t%q\n" $v}}{{end -}}
{{if .StdImports}}{{printf "\n"}}{{end -}}
{{range $v := .Imports}}{{printf "\t%q\n" $v}}{{end -}}
)

// {{.ComponentName}}Base ...
type {{.ComponentName}}Base struct{
	vecty.Core
}

// Render ...
func (c *{{.ComponentName}}Base) Render() vecty.ComponentOrHTML {
	return {{.Generated}}
}
{{range $method, $v := .Methods}}
// {{$method}} ...
func (c *{{$.ComponentName}}Base) {{$method}}(event *vecty.Event) {
	log.Print("{{$method}} not implement")
}
{{end}}
`))
