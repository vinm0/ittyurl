package web

import (
	"net/http"

	"github.com/vinm0/ittyurl/data"
)

// page contains the data to be processed within a template.
// Please use NewPage to initialize a new page instance.
// A title value is always required.
type page map[string]interface{}

// Get returns the value corresponding the named key.
// Returns nil if no such key exists.
func (p *page) Get(key string) interface{} {
	return (*p)[key]
}

// Add assigns a key-value pair to the Page.
// Add will overwrite any existing value with the same key.
func (p *page) Add(key string, val interface{}) {
	(*p)[key] = val
}

// Serve passes data to templates to include in an http response.
func (p *page) Serve(w http.ResponseWriter, templates ...string) {
	Render(p, w, templates...)
}

// NewPage returns a new Page instance populated with the provided title.
// A title is required for every page instance.
func NewPage(title string) *page {
	return &page{TITLE: title}
}

// Handles the root path and all undefined paths
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Unrecognized paths will be directed to custom 404
	if r.URL.Path == PATH_HOME {
		p := NewPage(TITLE_SITE)
		p.Serve(w, TMPL_BASE, TMPL_HOME)
		return
	}

	if path, found := data.RegisteredPath(r.URL.Path); found {
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	handleStaticPage(w, r)
	return
}

// Handles Post requests for new rows in the database
func handleNew(w http.ResponseWriter, r *http.Request) {
	if GetRequest(r) {
		handleStaticPage(w, r)
		return
	}

	category := r.URL.Query().Get("")

	switch category {
	case "url":
		NewUrl(w, r)

	case "user":
		NewUser(w, r)
	}
}

func handleSignin(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)

	if auth, _ := session.Values[SESSION_AUTH].(bool); auth {
		http.Redirect(w, r, PATH_PROFILE, http.StatusFound)
		return
	}

	if PostRequest(r) {
		Signin(w, r)
		return
	}

	p := NewPage(TITLE_SIGNIN)

	p.Serve(w, TMPL_BASE, TMPL_SIGNIN)
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
	default:
		title = TITLE_ERR
		templ = TMPL_ERR
	}

	p := NewPage(title)

	p.Serve(w, TMPL_BASE, templ)
}

func (s *Session) RedirectFlash(r *http.Request, w http.ResponseWriter, path, msg string) {
	s.AddFlash(msg)
	s.Save(r, w)
	http.Redirect(w, r, path, http.StatusFound)
}

func PostRequest(r *http.Request) bool {
	return r.Method == http.MethodPost
}

func GetRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}
