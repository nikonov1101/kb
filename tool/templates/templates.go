package templates

import (
	_ "embed"
	"html/template"
)

type Args struct {
	Title string
	Date  string

	Posts     []Post
	OrContent template.HTML
}

type Post struct {
	Title string
	Date  string

	Href string
}

//go:embed page.tmpl
var tmpl string

var _parsed *template.Template

func Load() *template.Template {
	if _parsed == nil {
		t, err := template.New("page").Parse(tmpl)
		if err != nil {
			panic(t)
		}
		_parsed = t
	}
	return _parsed
}
