package main

import (
	"log"

	"github.com/andyp1xe1/eden.nvim/eden/appview"
)

func main() {
	log.SetFlags(0)

	app := appview.MakeAppView(
		true,
		Name,
	)

	server := MakeServer(app.EvenHub())
	server.Serve()

	app.Run()
	app.Destroy()
}
