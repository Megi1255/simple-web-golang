package log

import (
	"context"
	"errors"
	"github.com/fluent/fluent-logger-golang/fluent"
	"log"
	"simple-web-golang/config"
	"time"
)

type Logger interface {
	Log(tag string, msg interface{}, time time.Time) error
}

func FromContext(c context.Context) Logger {
	val := c.Value(config.KeyLogger)
	return val.(Logger)
}

type Fluent struct {
	conf     *config.LoggerConfig
	fluent   *fluent.Fluent
	IsEnable bool
}

var (
	ErrNotAvailable = errors.New("[fluent] not available")
)

func New(cfg *config.LoggerConfig) *Fluent {
	logger, err := fluent.New(fluent.Config{
		FluentPort:  cfg.Port,
		FluentHost:  cfg.Host,
		Timeout:     cfg.ConnTimeout,
		BufferLimit: cfg.BufferLength,
		MaxRetry:    cfg.MaxConnTrial,
		TagPrefix:   cfg.TagPrefix,
	})
	if err != nil {
		log.Printf("[fluent] failed to connect: %s:%d", cfg.Host, cfg.Port)
		return &Fluent{
			conf:     cfg,
			fluent:   logger,
			IsEnable: false,
		}
	}
	return &Fluent{
		conf:     cfg,
		fluent:   logger,
		IsEnable: true,
	}
}

func (f *Fluent) Log(tag string, msg interface{}, time time.Time) error {
	if f.IsEnable {
		return f.fluent.PostWithTime(tag, time, msg)
	}
	return ErrNotAvailable
}
