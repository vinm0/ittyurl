package web

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/vinm0/ittyurl/data"
)

const (
	SESSION     = "session"
	SESSION_KEY = "SESSION_KEY"

	// Session keys

	SESSION_USR  = "usr"
	SESSION_AUTH = "authenticated"
)

var store *sessions.CookieStore

type Session struct{ *sessions.Session }

// Returns the the current session data.
func CurrentSession(r *http.Request) (session *Session) {
	s, _ := store.Get(r, SESSION)
	session = &Session{s}

	return session
}

// Initializes the session.
func SessionStart() {
	godotenv.Load(".env")
	store = sessions.NewCookieStore([]byte(os.Getenv(SESSION_KEY)))
}

// Clears the session of all values.
func (s *Session) Clear(w http.ResponseWriter, r *http.Request) {
	// s.Values = map[interface{}]interface{}{}
	s.Options.MaxAge = -1
	s.Save(r, w)
}

// Removes the key and it's corresponding value from the session map.
func (s *Session) Del(w http.ResponseWriter, r *http.Request, key string) {
	if _, ok := s.Values[key]; ok {
		delete(s.Values, key)
		s.Save(r, w)
	}
}

// Returns usr populated with User data from the session.
// Returns a public User, if no User data exists in the session.
func (s *Session) User() (usr *data.User) {
	usr, ok := s.Values[SESSION_USR].(*data.User)
	if !ok {
		usr = &data.User{
			UserID: data.USERID_PUBLIC,
			AcctID: data.ACCTTYPE_PUBL,
		}
	}
	return usr
}

// Adds a flash message to the session and redirects to the specified page.
// RedirectFlash is intended for informing the client of possible errors
// that have occurred during a POST request.
func (s *Session) RedirectFlash(r *http.Request, w http.ResponseWriter,
	path string, msg template.HTML) {

	s.Values["err"] = msg
	s.Save(r, w)
	http.Redirect(w, r, path, http.StatusFound)
}
