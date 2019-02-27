package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"simple-web-golang/api"
	"simple-web-golang/config"
)

var (
	ErrUnknownAPI = errors.New("Unknown API")
)

type Request struct {
	ApiName    string `json:"api_name" binding:"required"`
	RequestId  string `json:"request_id"`
	SessionKey string `json:"session_key"`
}

func (ctl *Controller) Gateway(c *gin.Context) {
	req, _ := ReqFrom(c)

	if err := invoke(&api.HandlerFunc{}, req.ApiName, c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func invoke(any interface{}, name string, args ...interface{}) error {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	val := reflect.ValueOf(any).MethodByName(name)
	if !val.IsValid() {
		return ErrUnknownAPI
	}
	val.Call(inputs)
	return nil
}

func ReqFrom(c *gin.Context) (req *Request, err error) {
	val := c.Value(config.KeyRequest)
	if val == nil {
		err = errors.New("not exist key: " + config.KeyRequest)
		return
	}
	req = val.(*Request)
	return
}
