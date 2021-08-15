package web

import (
	"net/http"
)

type Page map[string]interface{}

func (p *Page) Get(name string) interface{} {
	return (*p)[name]
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	p := &Page{
		TITLE: SITE_TITLE,
	}

	Render(w, p, TMPL_BASE, TMPL_HOME)
}

func handleSignout(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	session.Clear(w, r)

	http.Redirect(w, r, "/", http.StatusFound)
}
