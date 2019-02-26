package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"simple-web-golang/config"
	"time"
)

func TsFromContext(c context.Context) time.Time {
	val := c.Value(config.KeyTimestamp)
	return val.(time.Time)
}

func ProfilerFromContext(c *gin.Context) *Profiler {
	val := c.Value(config.KeyProfiler)
	return val.(*Profiler)
}
