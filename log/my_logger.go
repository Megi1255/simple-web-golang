package log

import (
	"log"
	"os"
	"time"
)

type MyLogger struct {
	Logger *log.Logger
}

func NewMyLogger(cfg *Config) *MyLogger {
	return &MyLogger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *MyLogger) Log(tag string, msg interface{}, time time.Time) error {
	log.Printf("[%s] %v", tag, msg)
	return nil
}
