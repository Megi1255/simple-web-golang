package log

import (
	"encoding/json"
	"log"
	"os"
	"simple-web-golang/config"
	"time"
)

type MyLogger struct {
	Logger *log.Logger
}

func NewMyLogger(cfg *config.LoggerConfig) *MyLogger {
	return &MyLogger{
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (l *MyLogger) Log(tag string, msg interface{}, time time.Time) error {
	entry, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[%s] %v", tag, msg)
	} else {
		log.Printf("[%s] %v", tag, string(entry))
	}
	return nil
}
