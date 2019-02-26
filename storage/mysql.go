package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"simple-web-golang/config"
)

type Rdb struct {
	Pool *sql.DB
}

func New(cfg *config.StorageConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true", cfg.User, cfg.Passwd, cfg.Host, cfg.Port, cfg.DbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("%s, %s", dsn, err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("%s, %s", dsn, err.Error())
	}
	return db
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
