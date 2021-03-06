package web

import (
	"html/template"
	"net/http"

	"github.com/vinm0/ittyurl/data"
)

// Creates a new Url object from POST form data.
// The url is saved to the session.
func NewUrl(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	usr := session.User()

	url, msg := usr.CreateUrl(r)
	if msg != "" {
		session.RedirectFlash(r, w, PATH_HOME, msg)
		return
	}

	session.Values["url"] = url
	session.Save(r, w)

	http.Redirect(w, r, PATH_HOME, http.StatusFound)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	usr := data.FormUser(r)

	if usr == nil {
		var msg template.HTML = "Something went wrong creating account. Try again"
		session.RedirectFlash(r, w, PATH_SIGNUP, msg)
		return
	}

	if msg := usr.CreateUser(); msg != "" {
		session.RedirectFlash(r, w, PATH_SIGNUP, msg)
		return
	}

	http.Redirect(w, r, PATH_SIGNIN, http.StatusFound)
}
