package web

import (
	"html/template"
	"log"
	"net/http"
)

/*
 * **********************************************
 * **********************************************
 * **************** Definitions *****************
 * **********************************************
 * **********************************************
 */

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

/*
 * **********************************************
 * **********************************************
 * ***************** Functions ******************
 * **********************************************
 * **********************************************
 */

// Render checks if a template exists for the page title.
// If a template exists, Render calls to execute the existing template.
// If a template does not exist, Render parses the specified templates,
// registers the template, and executes the new template.
func Render(w http.ResponseWriter, p *Page, templates ...string) {
	title, _ := p.Get(TITLE).(string)

	// Check if template exists in the template map.
	t, exists := tmap[title]
	if exists {
		t.Serve(w, p)
		return
	}

	templ := template.Must(template.ParseFiles(templates...))

	// Add new template to template map.
	t = &Template{templ}
	tmap[title] = t

	t.Serve(w, p)
}

// Serve executes the template provided by the receiver
func (t *Template) Serve(w http.ResponseWriter, p *Page) {
	err := t.template.Execute(w, p)
	if err != nil {
		log.Println(err.Error(), t.template)
	}
}

// Initializes the template map
func initTmap() {
	tmap = map[string]*Template{}
}
