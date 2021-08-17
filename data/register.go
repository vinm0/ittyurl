package data

import "encoding/gob"

func RegisterStructs() {
	gob.Register(User{})
	gob.Register(Url{})
}
