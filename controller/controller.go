package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct{}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func HandleError(c *gin.Context, code int, msg string) ErrorResponse {
	res := ErrorResponse{
		Code:  code,
		Error: msg,
	}
	c.JSON(http.StatusOK, res)
	return res
}
