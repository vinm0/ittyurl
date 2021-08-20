package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/vinm0/ittyurl/data"
)

/*
 * **********************************************
 * ************                    **************
 * ************    Page Content    **************
 * ************                    **************
 * **********************************************
 */

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

/*
 * **********************************************
 * **************                ****************
 * **************    Handlers    ****************
 * **************                ****************
 * **********************************************
 */

// Handles the root path and all undefined paths
func handleHome(w http.ResponseWriter, r *http.Request) {

	// Serve home page if the root path is the target
	if r.URL.Path == PATH_HOME {
		p := NewPage(TITLE_SITE)

		session := CurrentSession(r)
		if url, ok := session.Values["url"].(data.Url); ok {
			p.Add("url", url)
			session.Del(w, r, "url")
		}

		p.Serve(w, TMPL_BASE, TMPL_HOME)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/")

	if url, found := data.RegisteredPath(path); found {
		url.TrackVisit(r)
		fmt.Println("Redirecting to ", url.Source)
		http.Redirect(w, r, "https://"+url.Source, http.StatusSeeOther)
		return
	}

	handleStaticPage(w, r)
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

// GET: Serves the sign-in page if no client is signed in.
// Redirects to the client's profile page if a client is signed in.
//
// POST: Signs in client and redirects to client's profile page.
// If sign-in is unsuccessful, redirects to sign-in page.
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

// Removes the client's session data and redirects to the root page.
func handleSignout(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	session.Clear(w, r)

	http.Redirect(w, r, "/", http.StatusFound)
}

// Serves a static page or recognized in the domain,
// including a custom 404.
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

/*
 * **********************************************
 * ***************               ****************
 * ***************    Helpers    ****************
 * ***************               ****************
 * **********************************************
 */

// Adds a flash message to the session and redirects to the specified page.
// RedirectFlash is intended for informing the client of possible errors
// that have occurred during a POST request.
func (s *Session) RedirectFlash(r *http.Request, w http.ResponseWriter, path, msg string) {
	s.AddFlash(msg)
	s.Save(r, w)
	http.Redirect(w, r, path, http.StatusFound)
}

// Returns true if r contains a POST request
func PostRequest(r *http.Request) bool {
	return r.Method == http.MethodPost
}

// Returns true if r contains a GET request
func GetRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}
