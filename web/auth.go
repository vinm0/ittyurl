package web

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/vinm0/ittyurl/data"
	"golang.org/x/crypto/bcrypt"
)

const (
	// The minimum acceptable character length for a password field
	MIN_PWD_LEN = 8

	// The maximum acceptable character length for an email field
	MAX_EMAIL_CHAR = 320
)

// Signin adds client data to the session for valid logins.
// Successful logins are redirected to the profile page.
// Failed logins are redirected to the signin page.
func Signin(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)
	path := PATH_PROFILE

	// Retreive client data from database using login credentials
	usr, errMsg := UserSignin(r)
	auth := (errMsg == "")

	// Send failed login message to client as a flash
	if !auth {
		session.AddFlash(errMsg)
		path = PATH_SIGNIN
	}

	session.Values[SESSION_USR] = usr   // On fail, contains only submitted credentials
	session.Values[SESSION_AUTH] = auth // login succeeded or failed
	session.Save(r, w)

	http.Redirect(w, r, path, http.StatusFound)
}

// UserSignin validates the login credentials.
// Sucessful logins return usr populated with all client data
// from the database as a User instance and an empty errMsg.
// Failed logins return usr populated with submitted credentails
// and a helpful error message for the client.
func UserSignin(r *http.Request) (usr *data.User, errMsg string) {
	usr = &data.User{}
	usr.Pwd = strings.TrimSpace(r.PostFormValue("pwd"))
	usr.Email = strings.ToLower(strings.TrimSpace(r.PostFormValue("email")))

	// Check for valid fiends before querying the database
	if valid, msg := validFields(usr.Email, usr.Pwd); !valid {
		return usr, msg
	}

	dbUsr, found := data.FindUsr(usr.Email)

	if !found || !pwdMatch(dbUsr.Pwd, usr.Pwd) {
		return usr, "Incorrect email or password"
	}

	return dbUsr, ""
}

// pwdMatch compares the hashed password to the client submitted password.
// Returns true if the passwords match. Otherwise, returns false.
func pwdMatch(pwdDB, pwd string) bool {
	salt := os.Getenv("SALT")
	err := bcrypt.CompareHashAndPassword([]byte(pwdDB), []byte(salt+pwd))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

// validFields validates the login fields submitted by the client against field constraints.
// Returns true and an empty string if validation succeeds.
// Otherwise, returns false and a helpful error message for the client.
func validFields(email, pwd string) (valid bool, errMsg string) {
	if len(pwd) < MIN_PWD_LEN {
		return false, "Password must be " + strconv.Itoa(MIN_PWD_LEN) + " or more characters"
	}
	if !isEmail(email) {
		return false, "Please enter a valid email"
	}

	return true, ""
}

// isEmail validates the format of the email field.
func isEmail(email string) bool {
	match, _ := regexp.MatchString("^\\S{2,}@\\S{2,}\\.\\S{2,}$", email)
	validSize := len(email) <= MAX_EMAIL_CHAR

	return match && validSize
}
