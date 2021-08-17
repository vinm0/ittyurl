package web

import (
	"net/http"

	"github.com/vinm0/ittyurl/data"
)

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

}

func NewUser(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	usr := data.FormUser(r)

	if usr == nil {
		msg := "Something went wrong creating account. Try again"
		session.RedirectFlash(r, w, PATH_SIGNUP, msg)
		return
	}

	if msg := usr.CreateUser(); msg != "" {
		session.RedirectFlash(r, w, PATH_SIGNUP, msg)
		return
	}

	http.Redirect(w, r, PATH_SIGNIN, http.StatusFound)
}
