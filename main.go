package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-web-golang/cache"
	"simple-web-golang/config"
	"simple-web-golang/controller"
	"simple-web-golang/log"
	"simple-web-golang/storage"
	"simple-web-golang/util"
)

func main() {
	cfgName := flag.String("c", "gin", "config file name")
	flag.Parse()

	// Loading
	cfg, err := config.Load(*cfgName)
	if err != nil {
		panic(err)
	}
	ctl := &controller.Controller{
		DB:    storage.NewMongoClient(cfg.Db),
		Cache: cache.New(cfg.Cache),
		Logg:  log.New(cfg.Logger),
		//Logg:log.NewMyLogger(cfg.Logger),
		Cfg: cfg,
	}
	ctl.SetProfileHandler(func(key string, perf *util.Performance) error {
		return ctl.Logg.Log("profiler", map[string]interface{}{
			"key":  key,
			"perf": perf.ElapsedTime.Seconds(),
		}, perf.StartTime)
	})
	ctl.EnableProfiler(true)

	// Routing
	var r *gin.Engine
	r = gin.New()
	r.Use(ctl.Setup())
	r.Use(ctl.BodyDump(WriteBodyLog))
	r.Use(ctl.APIPerformance())

	r.POST("/", RootHandler)
	r.POST("/gateway", ctl.Gateway)

	// Start
	r.Run(fmt.Sprintf(":%d", cfg.Port))
}

func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to gin example!",
	})
}

func WriteBodyLog(c *gin.Context, reqBody, resBody []byte) {
	var req map[string]interface{}
	var res map[string]interface{}

	prof := util.ProfilerFromContext(c)
	prof.Start("middleware.bodydump")
	if err := json.Unmarshal(reqBody, &req); err != nil {
		return
	}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return
	}
	ts := util.TsFromContext(c)
	logg := log.FromContext(c)
	if err := logg.Log("game.request", req, ts); err != nil {
		fmt.Printf("failed to write log: %v", err)
		return
	}
	if err := logg.Log("game.response", res, ts); err != nil {
		fmt.Printf("failed to write log: %v", err)
		return
	}
	prof.End("middleware.bodydump")
}
