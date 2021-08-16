package web

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/vinm0/ittyurl/data"
	"golang.org/x/crypto/bcrypt"
)

const (
	MIN_PWD_LEN = 8
)

func Signin(w http.ResponseWriter, r *http.Request) {
	session := CurrentSession(r)

	usr, errMsg := UserSignin(r)
	if errMsg != "" {
		session.AddFlash(errMsg)
		session.Values["usr"] = usr
		session.Save(r, w)
		http.Redirect(w, r, PATH_SIGNIN, http.StatusFound)
		return
	}

}

func UserSignin(r *http.Request) (usr *data.User, errMsg string) {
	usr = &data.User{}
	usr.Pwd = r.PostFormValue("pwd")
	usr.Email = r.PostFormValue("email")

	if valid, msg := validCredentials(usr.Email, usr.Pwd); !valid {
		return usr, msg
	}

	dbUsr, auth := DBAuth(usr.Email)

	if !auth || !pwdMatch(dbUsr.Pwd, usr.Pwd) {
		return usr, "Incorrect email or password"
	}

	return dbUsr, ""
}

func pwdMatch(pwdDB, pwd string) bool {
	salt := os.Getenv("SALT")
	err := bcrypt.CompareHashAndPassword([]byte(pwdDB), []byte(salt+pwd))
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func validCredentials(email, pwd string) (valid bool, errMsg string) {
	if len(pwd) < MIN_PWD_LEN {
		return false, "Password must be " + strconv.Itoa(MIN_PWD_LEN) + " or more characters"
	}
	if !IsEmail(email) {
		return false, "Please enter a valid email"
	}

	return true, ""
}

func IsEmail(email string) bool {
	match, _ := regexp.MatchString("\\S{2,}@\\S{2,}\\.\\S{2,}$", email)
	return match
}
