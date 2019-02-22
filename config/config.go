package config

import (
	"ginsample/cache"
	"ginsample/log"
	"ginsample/storage"
	"github.com/spf13/viper"
)

const (
	ServerModeDebug = true
	ServerPort      = 1323
)

type Config struct {
	Debug     bool
	Port      int
	Db        *storage.Config
	Redis     *cache.Config
	Logger    *log.Config
	StoreType string
}

func Load(content string) (*Config, error) {
	config := &Config{}
	viper.SetDefault("port", ServerPort)
	viper.SetDefault("debug", ServerModeDebug)
	viper.SetDefault("db", storage.DefaultConfig())
	viper.SetDefault("redis", cache.DefaultConfig())
	viper.SetDefault("logger", log.DefaultConfig())
	viper.SetDefault("storetype", StoreTypeInMem)

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
