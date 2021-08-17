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

func (url *Url) TrackVisit(r *http.Request) {
	// TODO: Record visit if Owner AcctType permits.
	switch acct := url.Owner.AcctID; {
	case acct >= ACCTTYPE_PAID:
		vis := extractVisit(r)
		vis.InsertVisit()
	case acct >= ACCTTYPE_FREE:
		IncrementVisits(url.Path)
	}
}

func (u *Url) DuplicateSource() (original *Url, duplicate bool) {
	original = UrlBySource(u.Source)
	return original, (original != nil)
}

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
