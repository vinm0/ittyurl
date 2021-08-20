package data

import (
	"net/http"
	"time"
)

type Visit struct {
	Path    string
	Visdate time.Time
	Geo     *time.Location
	IP      string
}

// TODO: Add logic to record the visit into the database.
func (vis *Visit) InsertVisit() {
}

// Extract visit data from the request header.
func extractVisit(r *http.Request) *Visit {
	return &Visit{
		Path:    r.URL.Path,
		Visdate: time.Now(),
		Geo:     time.Now().Location(),
		IP:      ipAddr(r),
	}
}
