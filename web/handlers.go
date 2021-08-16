package web

import (
	"net/http"
)

type Page map[string]interface{}

func (p *Page) Get(name string) interface{} {
	return (*p)[name]
}

func (p *Page) Add(name string, val interface{}) {
	(*p)[name] = val
}

func newPage(title string) *Page {
	return &Page{TITLE: title}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	p := newPage(TITLE_SITE)

	p.Render(w, TMPL_BASE, TMPL_HOME)
}

func handleSignin(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)

	if auth, _ := session.Values["authenticated"].(bool); auth {
		http.Redirect(w, r, PATH_PROFILE, http.StatusFound)
		return
	}

	if PostMethod(r) {
		Signin(w, r)
		return
	}

	p := newPage(TITLE_SIGNIN)

	p.Render(w, TMPL_BASE, TMPL_SIGNIN)
}

func handleSignout(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	session.Clear(w, r)

	http.Redirect(w, r, "/", http.StatusFound)
}

func handleStaticPage(w http.ResponseWriter, r *http.Request) {
	title := ""
	templ := ""

	switch r.URL.Path {
	case PATH_PRIVACY:
		title = TITLE_PRIVACY
		templ = TMPL_PRIVACY
	case PATH_TERMS:
		title = TITLE_TERMS
		templ = TMPL_TERMS
	}

	p := newPage(title)

	p.Render(w, TMPL_BASE, templ)
}

func PostMethod(r *http.Request) bool {
	return r.Method == http.MethodPost
}
