package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"simple-web-golang/cache"
	"simple-web-golang/config"
	"simple-web-golang/log"
	"simple-web-golang/util"
	"time"
)

type Controller struct {
	DB    *mongo.Client
	Cache cache.Cache
	Logg  log.Logger
	Cfg   *config.Config

	profileOn      bool
	profileHandler util.FlushHandler
}

func (ctl *Controller) SetProfileHandler(h util.FlushHandler) {
	ctl.profileHandler = h
}

func (ctl *Controller) EnableProfiler(b bool) {
	ctl.profileOn = b
}

func (ctl *Controller) Setup() gin.HandlerFunc {
	//ctl.Logg.Log("middleware.setup", *ctl.Cfg, time.Now())
	return func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		ts := time.Now()
		prof := util.NewProfiler()
		prof.Enable = ctl.profileOn
		prof.SetFlushHandler(ctl.profileHandler)

		c.Set(config.KeyRequest, &req)
		c.Set(config.KeyConfig, ctl.Cfg)
		c.Set(config.KeyStorage, ctl.DB)
		c.Set(config.KeyCache, ctl.Cache)
		c.Set(config.KeyLogger, ctl.Logg)
		c.Set(config.KeyProfiler, prof)
		c.Set(config.KeyTimestamp, ts)
		c.Next()
		if err := prof.Flush(); err != nil {
			ctl.Logg.Log("profiler", gin.H{"error": err.Error()}, ts)
		}
	}
}

func (ctl *Controller) APIPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		prof := util.ProfilerFromContext(c)
		req, _ := ReqFrom(c)
		prof.Start("api." + req.ApiName)
		c.Next()
		prof.End("api." + req.ApiName)
	}
}
