package data

import (
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	URLS_TABLE = "urls JOIN users ON (owner_id = user_id)"
)

type Url struct {
	Path        string
	Source      string
	DateCreated time.Time
	Owner       *User
	CreatorIP   net.IP
}

// TODO: Record visit if Owner AcctType permits.
//
// Calls to update database based on permissions of the owner's account type.
func (url *Url) TrackVisit(r *http.Request) {
	switch acct := url.Owner.AcctID; {
	case acct >= ACCTTYPE_PAID:
		vis := extractVisit(r)
		vis.InsertVisit()
	case acct >= ACCTTYPE_FREE:
		IncrementVisits(url.Path)
	}
}

// Returns the Url instance if it exists in the database.
//
// If Url does not exist, original is nil.
func (u *Url) DuplicateSource() (original *Url, duplicate bool) {
	original = UrlBySource(u.Source)
	return original, (original != nil)
}

// Returns the Url data extracted from the POST request.
//
// A URL instance is always returned.
func UrlFromForm(r *http.Request, usr *User) *Url {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), " ")[0])

	return &Url{
		Path:        RandomPath(),
		Source:      r.PostFormValue("source"),
		DateCreated: time.Now(),
		Owner:       usr,
		CreatorIP:   ip,
	}
}
