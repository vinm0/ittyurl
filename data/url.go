package data

import (
	"net/http"
	"regexp"
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
	// TODO: Validate Source url format (http[s]://), trim values
	src := strings.TrimSpace(r.PostFormValue("source"))

	return &Url{
		Path:        RandomPath(),
		Source:      formatUrl(src),
		DateCreated: time.Now(),
		Owner:       usr,
		CreatorIP:   ipAddr(r),
	}
}

// Returns an string representation of the ip address extracted from
// the http request.
func ipAddr(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = strings.Split(r.Header.Get("X-Forwarded-For"), ", ")[0]
	}
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

// Returns a url string with "https://" prepended.
// Returns the original stirng if it already contains the http[s] protocal
func formatUrl(source string) string {
	if match, _ := regexp.MatchString("https?:\\/\\/", source); match {
		return source
	}

	return "https://" + source
}
