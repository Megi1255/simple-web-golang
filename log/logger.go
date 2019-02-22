package log

import "time"

const (
	DefaultHost         = "127.0.0.1"
	DefaultPort         = 24224
	DefaultBufferLength = 1000
	DefaultMaxConnTrial = 5
	DefaultConnTimeout  = 5 * time.Second
	DefaultFlushTimeout = 500 * time.Millisecond
	DefaultTagPrefix    = ""
)

type Logger interface {
	Log(tag string, msg interface{}, time time.Time) error
}

type Config struct {
	Host         string
	Port         int
	BufferLength int
	MaxConnTrial int
	ConnTimeout  time.Duration
	FlushTimeout time.Duration
	TagPrefix    string
}

func DefaultConfig() *Config {
	return &Config{
		Host:         DefaultHost,
		Port:         DefaultPort,
		BufferLength: DefaultBufferLength,
		MaxConnTrial: DefaultMaxConnTrial,
		ConnTimeout:  DefaultConnTimeout,
		FlushTimeout: DefaultFlushTimeout,
		TagPrefix:    DefaultTagPrefix,
	}
}
