package main

import (
	"code.google.com/p/go.mobile/app"

	_ "code.google.com/p/go.mobile/bind/java"
	_ "github.com/howeyc/spipedmobile/go_spiped"
)

func main() {
	app.Run(app.Callbacks{})
}
