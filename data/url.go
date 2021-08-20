package data

import (
	"net/http"
	"strings"
	"time"
)

const (
	URLS_TABLE = "urls JOIN usrs ON (owner_id = usr_id)"
)

type Url struct {
	Path        string
	Source      string
	DateCreated time.Time
	Owner       *User
	CreatorIP   string
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
func UrlFromPost(r *http.Request, usr *User) *Url {
	return &Url{
		Path:        RandomPath(),
		Source:      r.PostFormValue("source"),
		DateCreated: time.Now(),
		Owner:       usr,
		CreatorIP:   ipAddr(r),
	}
}

func ipAddr(r *http.Request) string {
	return strings.Split(r.Header.Get("X-Forwarded-For"), ", ")[0]
}
