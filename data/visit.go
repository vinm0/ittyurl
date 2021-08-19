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

// TODO: Add logic to record the visit into the database.
func (vis *Visit) InsertVisit() {
}

// Extract visit data from the request header.
func extractVisit(r *http.Request) *Visit {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), " ")[0])

	return &Visit{
		Path:    r.URL.Path,
		Visdate: time.Now(),
		Geo:     time.Now().Location(),
		IP:      ip,
	}
}
