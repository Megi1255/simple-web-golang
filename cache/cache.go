package cache

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

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Exists(key string) (bool, error)
	UpdateRank(key string, point int, uid int64) error
	Rank(key string, uid int64) (int64, error)
	Flush() error
	SetMap(key string, value map[string]interface{}) error
	GetMap(key string) (map[string]interface{}, error)
}

type Config struct {
	Host   string
	Port   int
	Db     int
	Expire int
	Prefix string
}

func DefaultConfig() *Config {
	return &Config{
		Host:   DefaultRedisHost,
		Port:   DefaultRedisPort,
		Prefix: DefaultRedisKeyPrefix,
		Db:     DefaultRedisDb,
		Expire: DefaultRedisKeyExpire,
	}
}
