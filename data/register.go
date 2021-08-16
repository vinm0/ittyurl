package data

import "encoding/gob"

func RegisterStructs() {
	gob.Register(User{})
}
