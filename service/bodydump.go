package service

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type BodyDumpHandler func(*gin.Context, []byte, []byte)
type bodyDumpResponseWriter struct {
	gin.ResponseWriter
	ResBody *bytes.Buffer
}

func BodyDump(handler BodyDumpHandler) gin.HandlerFunc {
	if handler == nil {
		panic("bodydump handler required")
	}

	return func(c *gin.Context) {
		var reqBody []byte
		if cb, ok := c.Get(gin.BodyBytesKey); ok {
			if cbb, ok := cb.([]byte); ok {
				reqBody = cbb
			}
		}
		if reqBody == nil {
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Set(gin.BodyBytesKey, body)
			reqBody = body
		}
		bdw := &bodyDumpResponseWriter{
			ResponseWriter: c.Writer,
			ResBody:        bytes.NewBuffer([]byte{}),
		}
		c.Writer = bdw
		c.Next()
		handler(c, reqBody, bdw.ResBody.Bytes())
	}
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	w.ResBody.Write(b)
	return w.ResponseWriter.Write(b)
}
