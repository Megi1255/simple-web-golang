package main

import (
	"ginsample/service"
)

func main() {
	app := service.New("gin")
	app.Run()
}
