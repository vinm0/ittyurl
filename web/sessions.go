package web

import (
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

func CurrentSession(r *http.Request) (session *Session) {
	s, _ := store.Get(r, SESSION)
	session = &Session{s}

	return session
}

func SessionStart() {
	godotenv.Load(".env")
	store = sessions.NewCookieStore([]byte(os.Getenv(SESSION_KEY)))
}

func (s *Session) Clear(w http.ResponseWriter, r *http.Request) {
	s.Values = map[interface{}]interface{}{}
	s.Save(r, w)
}

func (s *Session) User() (usr *data.User) {
	usr, ok := s.Values[SESSION_USR].(*data.User)
	if !ok {
		usr = &data.User{UserID: data.USERID_PUB}
	}
	return usr
}
