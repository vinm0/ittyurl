package data

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
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

// Inserts a User record in to the database.
// Returns an error message if an issues occurs.
func (usr *User) CreateUser() (errMsg string) {
	if err := createUser(usr); err != nil {
		errMsg = "There was an issue reaching the database"
	}

	return errMsg
}

// Creates a Url from the POST form values, and returns the Url.
//
// If the site's policies permit, the Url is inserted into the database.
// If the site's policies do not permit, an errMsg is returned with the new Url.
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

	if err := insertUrl(url); err != nil {
		errMsg = "There was a problem connecting to the database."
	}

	return url, errMsg
}

// TODO: Query validate limits.
//
// Returns true if the user has exceeded the account type's url limits.
func (usr *User) exceedUrlLmit(ip net.IP) bool {
	return false
}

// Extracts the User data from the POST form values.
//
// A new User instance is always returned.
func FormUser(r *http.Request) (usr *User) {
	id, _ := strconv.Atoi(r.PostFormValue("acctid"))

	usr = &User{
		Fname:    r.PostFormValue("fname"),
		Lname:    r.PostFormValue("lname"),
		Email:    r.PostFormValue("email"),
		Usrname:  r.PostFormValue("usrname"),
		Pwd:      r.PostFormValue("pwd"),
		Joindate: time.Now(),
		AcctID:   id,
	}
	return usr
}
