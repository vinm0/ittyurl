package web

import (
	"net/http"
)

type Page map[string]interface{}

func (p *Page) Get(name string) interface{} {
	return (*p)[name]
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

func PostMethod(r *http.Request) bool {
	return r.Method == http.MethodPost
}
