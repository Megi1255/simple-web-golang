package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"simple-web-golang/cache"
	"simple-web-golang/config"
	"simple-web-golang/log"
	"simple-web-golang/storage"
	"time"
)

const (
	KeyPrefix    = "ginsample::"
	KeyConfig    = KeyPrefix + "config"
	KeyStorage   = KeyPrefix + "storage"
	KeyCache     = KeyPrefix + "cache"
	KeyLogger    = KeyPrefix + "logger"
	KeyTimestamp = KeyPrefix + "timestamp"
)

type Controller struct{}

func Setup(cfg *config.Config) gin.HandlerFunc {
	var db *sql.DB
	var redis cache.Cache
	var logg log.Logger
	db = storage.New(cfg.Db)
	redis = cache.New(cfg.Cache)
	//logg = log.New(cfg.Logger)
	logg = log.NewMyLogger(cfg.Logger)

	return func(c *gin.Context) {
		c.Set(KeyConfig, cfg)
		c.Set(KeyStorage, db)
		c.Set(KeyCache, redis)
		c.Set(KeyLogger, logg)
		c.Set(KeyTimestamp, time.Now())
		c.Next()
	}
}

func ConfFrom(c *gin.Context) (conf *config.Config, err error) {
	val, ok := c.Get(KeyConfig)
	if !ok {
		err = errors.New("not exist key: " + KeyConfig)
		return
	}
	conf = val.(*config.Config)
	return
}

func DBFrom(c *gin.Context) (db *sql.DB, err error) {
	val, ok := c.Get(KeyStorage)
	if !ok {
		err = errors.New("not exist key: " + KeyStorage)
		return
	}
	db = val.(*sql.DB)
	return
}

func CacheFrom(c *gin.Context) (cac cache.Cache, err error) {
	val, ok := c.Get(KeyCache)
	if !ok {
		err = errors.New("not exist key: " + KeyCache)
		return
	}
	cac = val.(cache.Cache)
	return
}

func LoggerFrom(c *gin.Context) (logger log.Logger, err error) {
	val, ok := c.Get(KeyLogger)
	if !ok {
		err = errors.New("not exist key: " + KeyLogger)
		return
	}
	logger = val.(log.Logger)
	return
}

func TsFrom(c *gin.Context) (ts time.Time, err error) {
	val, ok := c.Get(KeyTimestamp)
	if !ok {
		err = errors.New("not exist key: " + KeyTimestamp)
		return
	}
	ts = val.(time.Time)
	return
}

func WriteBodyLog(c *gin.Context, reqBody, resBody []byte) {
	var req map[string]interface{}
	var res map[string]interface{}

	if err := json.Unmarshal(reqBody, &req); err != nil {
		return
	}
	if err := json.Unmarshal(resBody, &res); err != nil {
		return
	}
	ts, _ := TsFrom(c)
	logg, _ := LoggerFrom(c)
	if err := logg.Log("game.request", req, ts); err != nil {
		return
	}
	if err := logg.Log("game.response", res, ts); err != nil {
		return
	}
}
