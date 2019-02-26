package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"reflect"
	"simple-web-golang/cache"
	"simple-web-golang/config"
	"simple-web-golang/log"
	"simple-web-golang/storage"
	"simple-web-golang/util"
	"time"
)

var (
	ErrUnknownAPI = errors.New("Unknown API")
)

type Request struct {
	ApiName    string `json:"api_name" binding:"required"`
	RequestId  string `json:"request_id"`
	SessionKey string `json:"session_key"`
}

func Gateway(c *gin.Context) {
	req, _ := ReqFrom(c)

	if err := Invoke(&Controller{}, req.ApiName, c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func Invoke(any interface{}, name string, args ...interface{}) error {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	val := reflect.ValueOf(any).MethodByName(name)
	if !val.IsValid() {
		return ErrUnknownAPI
	}
	val.Call(inputs)
	return nil
}

func ReqFrom(c *gin.Context) (req *Request, err error) {
	val := c.Value(config.KeyRequest)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyRequest)
		return
	}
	req = val.(*Request)
	return
}

func Setup(cfg *config.Config) gin.HandlerFunc {
	ts := time.Now()
	//db = storage.New(cfg.Db)
	db := storage.NewMongoClient(cfg.Db)
	redis := cache.New(cfg.Cache)
	//logg = log.New(cfg.Logger)
	logg := log.NewMyLogger(cfg.Logger)
	logg.Log("middleware.setup", *cfg, ts)

	return func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		prof := util.NewProfiler()
		prof.SetFlushHandler(func(key string, perf util.Performance) error {
			fmt.Printf("%s %v", key, perf)
			return logg.Log("profiler", map[string]interface{}{
				"key":  key,
				"perf": perf.ElapsedTime.Nanoseconds(),
			}, ts)
		})

		c.Set(config.KeyRequest, &req)
		c.Set(config.KeyConfig, cfg)
		c.Set(config.KeyStorage, db)
		c.Set(config.KeyCache, redis)
		c.Set(config.KeyLogger, logg)
		c.Set(config.KeyProfiler, prof)
		c.Set(config.KeyTimestamp, ts)
		c.Next()
		if err := prof.Flush(); err != nil {
			logg.Log("profiler", gin.H{"error": err.Error()}, ts)
		}
	}
}

func APIPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		prof := util.ProfilerFromContext(c)
		req, _ := ReqFrom(c)
		prof.Start("api." + req.ApiName)
		c.Next()
		prof.End("api." + req.ApiName)
	}
}
