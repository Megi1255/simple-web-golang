package config

import (
	"errors"
)

const (
	DefaultRedisHost      = "127.0.0.1"
	DefaultRedisPort      = 6379
	DefaultRedisKeyPrefix = "Gin::"
	DefaultRedisDb        = 0
	DefaultRedisKeyExpire = 60
)

var (
	ErrUserNotExist     = errors.New("User does not exist")
	ErrUserAlreadyExist = errors.New("User already exist")
	ErrCacheMissingKey  = errors.New("cache missing key")
)

type CacheConfig struct {
	Host   string
	Port   int
	Db     int
	Expire int
	Prefix string
}

func CacheDefaultConfig() *CacheConfig {
	return &CacheConfig{
		Host:   DefaultRedisHost,
		Port:   DefaultRedisPort,
		Prefix: DefaultRedisKeyPrefix,
		Db:     DefaultRedisDb,
		Expire: DefaultRedisKeyExpire,
	}
}
