package log

import (
	"errors"
	"github.com/fluent/fluent-logger-golang/fluent"
	"log"
	"time"
)

type Fluent struct {
	conf     *Config
	fluent   *fluent.Fluent
	IsEnable bool
}

var (
	ErrNotAvailable = errors.New("[fluent] not available")
)

func New(cfg *Config) *Fluent {
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
