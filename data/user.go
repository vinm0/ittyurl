package data

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	USERID_PUBLIC = 0

	ACCTTYPE_PUBL = 0
	ACCTTYPE_FREE = 1000
	ACCTTYPE_PAID = 2000
	ACCTTYPE_PREM = 3000

	TIME_FORMAT = "2006-01-02 15:04"
)

type User struct {
	UserID   int
	Fname    string
	Lname    string
	Email    string
	Usrname  string
	Pwd      string
	Joindate time.Time
	AcctID   int
}

func (usr *User) CreateUser() (errMsg string) {
	if err := createUser(usr); err != nil {
		errMsg = "There was an issue reaching the database"
	}

	return errMsg
}

func (usr *User) CreateUrl(r *http.Request) (url *Url, errMsg string) {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), " ")[0])

	fmt.Println("Raw Header IP:", r.Header.Get("X-Forwarded-For"))
	fmt.Println("Formatted IP:", ip)

	if usr.exceedUrlLmit(ip) {
		return nil, "Exceeded url creation limits"
	}

	url = UrlFromForm(r, usr)
	if u, exists := url.DuplicateSource(); exists {
		return u, ""
	}

	if err := createUrl(url); err != nil {
		errMsg = "There was a problem connecting to the database."
	}

	return url, errMsg
}

func (usr *User) exceedUrlLmit(ip net.IP) bool {
	return false
}

func FormUser(r *http.Request) (usr *User) {
	return nil
}
