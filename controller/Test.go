package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func (h *Controller) Test(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ts, _ := TsFrom(c)
	logger, err := LoggerFrom(c)
	req["logging"] = "false"
	req["timestamp"] = ts.String()
	if err == nil {
		if err := logger.Log("test", req, ts); err == nil {
			req["logging"] = "true"
		} else {
			req["logging"] = err.Error()
		}
	}

	c.JSON(http.StatusOK, req)
	return
}
