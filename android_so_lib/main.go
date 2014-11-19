package main

import (
	"golang.org/x/mobile/app"

	_ "golang.org/x/mobile/bind/java"
	_ "github.com/howeyc/spipedmobile/go_spiped"
)

func main() {
	app.Run(app.Callbacks{})
}
