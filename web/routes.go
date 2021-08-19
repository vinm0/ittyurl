package web

import (
	"log"
	"net/http"
)

const (
	PORT = ":8080"

	// Domain Paths
	PATH_HOME    = "/"
	PATH_SIGNIN  = "/signin/"
	PATH_SIGNUP  = "/signup/"
	PATH_SIGNOUT = "/signout/"
	PATH_PROFILE = "/profile/"
	PATH_PRIVACY = "/privacy/"
	PATH_TERMS   = "/terms/"
	PATH_NEW     = "/new/"
)

func LaunchRouter() {
	http.HandleFunc(PATH_SIGNIN, handleSignin)
	http.HandleFunc(PATH_SIGNOUT, handleSignout)
	http.HandleFunc(PATH_PRIVACY, handleStaticPage)
	http.HandleFunc(PATH_TERMS, handleStaticPage)

	http.HandleFunc(PATH_NEW, handleNew)

	http.HandleFunc(PATH_HOME, handleHome)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))),
	)

	log.Fatal(http.ListenAndServe(PORT, nil))
}
