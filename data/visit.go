package data

import (
	"net"
	"net/http"
	"strings"
	"time"
)

type Visit struct {
	Path    string
	Visdate time.Time
	Geo     *time.Location
	IP      net.IP
}

func (vis *Visit) InsertVisit() {
	//
}

func extractVisit(r *http.Request) *Visit {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), " ")[0])

	return &Visit{
		Path:    r.URL.Path,
		Visdate: time.Now(),
		Geo:     time.Now().Location(),
		IP:      ip,
	}
}
