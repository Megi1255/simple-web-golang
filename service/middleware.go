package service

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"simple-web-golang/cache"
	"simple-web-golang/config"
	"simple-web-golang/log"
	"simple-web-golang/storage"
	"simple-web-golang/util"
	"time"
)

func Setup(cfg *config.Config) gin.HandlerFunc {
	var db *sql.DB
	var redis cache.Cache
	var logg log.Logger
	db = storage.New(cfg.Db)
	redis = cache.New(cfg.Cache)
	//logg = log.New(cfg.Logger)
	logg = log.NewMyLogger(cfg.Logger)

	return func(c *gin.Context) {
		c.Set(config.KeyConfig, cfg)
		c.Set(config.KeyStorage, db)
		c.Set(config.KeyCache, redis)
		c.Set(config.KeyLogger, logg)
		c.Set(config.KeyTimestamp, time.Now())
		c.Next()
	}
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
	ts, _ := util.TsFrom(c)
	logg, _ := util.LoggerFrom(c)
	if err := logg.Log("game.request", req, ts); err != nil {
		return
	}
	if err := logg.Log("game.response", res, ts); err != nil {
		return
	}
}
