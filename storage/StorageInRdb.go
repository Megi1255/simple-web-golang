package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Rdb struct {
	Pool *sql.DB
}

func New(cfg *Config) *sql.DB {
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
