package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-web-golang/config"
	"simple-web-golang/controller"
)

func main() {
	cfgName := flag.String("c", "gin", "config file name")
	flag.Parse()

	cfg, err := config.Load(*cfgName)
	if err != nil {
		panic(err)
	}
	var r *gin.Engine
	r = gin.New()

	r.Use(controller.Setup(cfg))
	r.Use(controller.BodyDump(controller.WriteBodyLog))
	r.Use(controller.APIPerformance())
	r.POST("/", RootHandler)
	r.POST("/gateway", controller.Gateway)
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}

func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to gin example!",
	})
}
