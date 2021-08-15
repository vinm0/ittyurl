package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	TMPL_DIR  = "templates/"
	TMPL_HOME = TMPL_DIR + "index.html"
	TMPL_BASE = TMPL_DIR + "base.html"

	SITE_TITLE = "IttyURL"
	TITLE      = "title"
)

var (
	tmap map[string]*Template
)

type Template struct {
	template *template.Template
}

func Render(w http.ResponseWriter, p *Page, templates ...string) {
	title, _ := p.Get(TITLE).(string)

	fmt.Println("Received title", title)

	t, exists := tmap[title]
	fmt.Println(title, "exists:", exists)
	if exists {
		fmt.Println(p)
		t.Serve(w, p)
		return
	}

	fmt.Println("Creating new template from", templates)

	templ := template.Must(template.ParseFiles(templates...))

	fmt.Println("Adding template to map")

	t = &Template{templ}
	tmap[title] = t

	fmt.Println("Executing Template")
	err := templ.Execute(w, p)
	if err != nil {
		log.Println(err.Error(), "unable to execute template", templ)
	}
}

func (t *Template) Serve(w http.ResponseWriter, p *Page) {
	t.template.Execute(w, p)
}

func initTmap() {
	tmap = map[string]*Template{}
}
