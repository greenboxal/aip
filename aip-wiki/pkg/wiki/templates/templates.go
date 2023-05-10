package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var content embed.FS
var templates *template.Template

func init() {
	t, err := template.ParseFS(content, "*.tmpl.html")

	if err != nil {
		panic(err)
	}

	templates = t
}

func Content() embed.FS {
	return content
}

func Templates() *template.Template {
	return templates
}
