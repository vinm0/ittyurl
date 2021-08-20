package data

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	USERID_PUBLIC = 0

	ACCTTYPE_PUBL = 0
	ACCTTYPE_FREE = 100
	ACCTTYPE_PAID = 200
	ACCTTYPE_PREM = 300

	MAX_URLS_PAID = 20
	MAX_URLS_FREE = 10
	MAX_URLS_PUBL = 5

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
func (usr *User) CreateUser() (errMsg template.HTML) {
	if err := createUser(usr); err != nil {
		errMsg = "There was an issue reaching the database"
	}

	return errMsg
}

// Creates a Url from the POST form values, and returns the Url.
//
// If the site's policies permit, the Url is inserted into the database.
// If the site's policies do not permit, an errMsg is returned with the new Url.
//
// A Url pointer is always returned
func (usr *User) CreateUrl(r *http.Request) (url *Url, errMsg template.HTML) {
	ip := ipAddr(r)

	fmt.Println("Raw Header IP:", r.Header.Get("X-Forwarded-For"))
	fmt.Println("Formatted IP:", ip)

	if usr.exceedUrlLmit(ip) {
		return nil,
			"Public url limits exceeded.<br><a href='/signin'><strong>Sign-in</strong></a> for more urls"
	}

	// format the "source" form value to valid url address
	r.ParseForm()
	r.PostForm.Set("source", formatUrl(r.PostFormValue("source")))

	// If Url already exists in the database, return the exising Url.
	url = UrlBySource(r.PostFormValue("source"))
	if url != nil {
		fmt.Println("previous url found.", url)
		return url, ""
	}

	// Extract Url from POST form
	url = UrlFromPost(r, usr)

	// Communicate error with client if database error occurs.
	if err := insertUrl(url); err != nil {
		errMsg = "There was a problem connecting to the database."
	}

	return url, errMsg
}

// TODO: Query validate limits.
//
// Returns true if the user has exceeded the account type's url limits.
func (usr *User) exceedUrlLmit(ip string) bool {
	if usr.AcctID >= ACCTTYPE_PREM {
		return false
	}
	if usr.AcctID == ACCTTYPE_PUBL {
		count := CountUrlsByIP(ip)
		fmt.Println("Urls Created:", count)
		return count >= MAX_URLS_PUBL
	}
	count := CountUrlsByUserID(usr.UserID)
	switch usr.AcctID {
	case ACCTTYPE_PAID:
		return count >= MAX_URLS_PAID
	case ACCTTYPE_FREE:
		return count >= MAX_URLS_FREE
	}
	return false
}

// Extracts the User data from the POST form values.
//
// A new User instance is always returned.
func FormUser(r *http.Request) (usr *User) {
	id, _ := strconv.Atoi(r.PostFormValue("acctid"))

	usr = &User{
		Fname:    strings.TrimSpace(r.PostFormValue("fname")),
		Lname:    strings.TrimSpace(r.PostFormValue("lname")),
		Email:    strings.TrimSpace(r.PostFormValue("email")),
		Usrname:  strings.TrimSpace(r.PostFormValue("usrname")),
		Pwd:      strings.TrimSpace(r.PostFormValue("pwd")),
		Joindate: time.Now(),
		AcctID:   id,
	}
	return usr
}
