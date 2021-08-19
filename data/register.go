package data

import "encoding/gob"

// Encodes the application's structs
func RegisterStructs() {
	gob.Register(User{})
	gob.Register(Url{})
}
