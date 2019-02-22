package main

import (
	"simple-web-golang/service"
)

func main() {
	app := service.New("gin")
	app.Run()
}
