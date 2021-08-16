package main

import (
	"github.com/vinm0/ittyurl/data"
	"github.com/vinm0/ittyurl/web"
)

func main() {
	web.SessionStart()
	data.RegisterStructs()

	web.Launch()
}
