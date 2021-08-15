package web

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

const (
	SESSION = "session"
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
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func (s *Session) Clear(w http.ResponseWriter, r *http.Request) {
	s.Values = map[interface{}]interface{}{}
	s.Save(r, w)
}
