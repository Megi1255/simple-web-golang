package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"reflect"
	"simple-web-golang/config"
	"simple-web-golang/controller"
	"time"
)

var (
	ErrUnknownAPI = errors.New("Unknown API")
)

type App struct {
	Config *config.Config
}

func New(cname string) *App {
	cfg, err := config.Load(cname)
	if err != nil {
		panic(err)
	}
	return &App{Config: cfg}
}

func (a *App) Run() {
	var r *gin.Engine
	r = gin.New()
	r.Use(controller.Setup(a.Config))
	r.Use(BodyDump(controller.WriteBodyLog))
	r.POST("/", RootHandler)
	r.POST("/gateway", Gateway)
	log.Fatal(r.Run(fmt.Sprintf(":%d", a.Config.Port)))
}

type Request struct {
	ApiName    string `json:"api_name" binding:"required"`
	RequestId  string `json:"request_id"`
	SessionKey string `json:"session_key"`
}

func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to gin example!",
	})
}

func Gateway(c *gin.Context) {
	start := time.Now()

	var req Request
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := Invoke(&controller.Controller{}, req.ApiName, c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	elapsed := time.Since(start)
	logg, _ := controller.LoggerFrom(c)
	ts, _ := controller.TsFrom(c)
	logg.Log("performance."+req.ApiName, map[string]interface{}{
		"api_name": req.ApiName,
		"time":     elapsed.Nanoseconds(),
	}, ts)
}

func Invoke(any interface{}, name string, args ...interface{}) error {
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
