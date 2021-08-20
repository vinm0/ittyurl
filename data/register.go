package data

import (
	"encoding/gob"
	"html/template"
)

// Encodes the application's structs
func RegisterStructs() {
	gob.Register(User{})
	gob.Register(Url{})
	var h template.HTML
	gob.Register(h)
}
