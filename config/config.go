package config

import (
	"github.com/spf13/viper"
	"simple-web-golang/cache"
	"simple-web-golang/log"
	"simple-web-golang/storage"
)

const (
	ServerModeDebug = true
	ServerPort      = 1323

	KeyPrefix    = "ginsample::"
	KeyConfig    = KeyPrefix + "config"
	KeyStorage   = KeyPrefix + "storage"
	KeyCache     = KeyPrefix + "cache"
	KeyLogger    = KeyPrefix + "logger"
	KeyTimestamp = KeyPrefix + "timestamp"
)

type Config struct {
	Debug     bool
	Port      int
	Db        *storage.Config
	Cache     *cache.Config
	Logger    *log.Config
	StoreType string
}

func Load(content string) (*Config, error) {
	config := &Config{}
	viper.SetDefault("port", ServerPort)
	viper.SetDefault("debug", ServerModeDebug)
	viper.SetDefault("db", storage.DefaultConfig())
	viper.SetDefault("cache", cache.DefaultConfig())
	viper.SetDefault("logger", log.DefaultConfig())

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
