package config

import (
	"context"
	"github.com/spf13/viper"
)

const (
	ServerModeDebug = true
	ServerPort      = 1323

	KeyPrefix    = "ginsample::"
	KeyRequest   = KeyPrefix + "request"
	KeyConfig    = KeyPrefix + "config"
	KeyStorage   = KeyPrefix + "storage"
	KeyCache     = KeyPrefix + "cache"
	KeyLogger    = KeyPrefix + "logger"
	KeyProfiler  = KeyPrefix + "profiler"
	KeyTimestamp = KeyPrefix + "timestamp"

	CODE_OK       = 200
	CODE_DB_ERROR = 500
)

type Config struct {
	Debug     bool
	Port      int
	Db        *StorageConfig
	Cache     *CacheConfig
	Logger    *LoggerConfig
	StoreType string
}

func Load(content string) (*Config, error) {
	config := &Config{}
	viper.SetDefault("port", ServerPort)
	viper.SetDefault("debug", ServerModeDebug)
	viper.SetDefault("db", StorageDefaultConfig())
	viper.SetDefault("cache", CacheDefaultConfig())
	viper.SetDefault("logger", LoggerDefaultConfig())

	var err error
	viper.SetConfigName(content)
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func FromContext(c context.Context) *Config {
	val := c.Value(KeyConfig)
	return val.(*Config)
}
