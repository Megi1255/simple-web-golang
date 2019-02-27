package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc struct{}

type ErrorResponse struct {
	ApiName string `json:"api_name"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
}

func HandleError(c *gin.Context, apiname string, code int, msg string) ErrorResponse {
	res := ErrorResponse{
		Code:  code,
		Error: msg,
	}
	c.JSON(http.StatusOK, res)
	return res
}
