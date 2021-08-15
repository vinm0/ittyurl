package web

import (
	"log"
	"net/http"
)

const (
	PORT = ":8080"
)

func Launch() {
	initTmap()
	SessionStart()

	http.HandleFunc("/signout", handleSignout)
	http.HandleFunc("/", handleHome)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))),
	)

	log.Fatal(http.ListenAndServe(PORT, nil))
}
