package util

import (
	"context"
	"database/sql"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"simple-web-golang/cache"
	"simple-web-golang/config"
	"simple-web-golang/log"
	"time"
)

func ConfFrom(c context.Context) (conf *config.Config, err error) {
	val := c.Value(config.KeyConfig)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyConfig)
		return
	}
	conf = val.(*config.Config)
	return
}

func DBFrom(c context.Context) (db *sql.DB, err error) {
	val := c.Value(config.KeyStorage)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyStorage)
		return
	}
	db = val.(*sql.DB)
	return
}

func CacheFrom(c context.Context) (cac cache.Cache, err error) {
	val := c.Value(config.KeyCache)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyCache)
		return
	}
	cac = val.(cache.Cache)
	return
}

func LoggerFrom(c context.Context) (logger log.Logger, err error) {
	val := c.Value(config.KeyLogger)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyLogger)
		return
	}
	logger = val.(log.Logger)
	return
}

func TsFrom(c context.Context) (ts time.Time, err error) {
	val := c.Value(config.KeyTimestamp)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyTimestamp)
		return
	}
	ts = val.(time.Time)
	return
}

func MongoFrom(c context.Context) (client *mongo.Client, err error) {
	val := c.Value(config.KeyStorage)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyTimestamp)
		return
	}
	client = val.(*mongo.Client)
	return
}
