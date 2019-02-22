package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"simple-web-golang/model"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Result int        `json:"result"`
	Error  string     `json:"error"`
	User   model.User `json:"user"`
}

func (h *Controller) Login(c *gin.Context) {
	var req LoginRequest
	var res LoginResponse
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ts, _ := TsFrom(c)
	cac, _ := CacheFrom(c)
	db, _ := DBFrom(c)
	logg, _ := LoggerFrom(c)

	var u model.User
	var err error
	if u, err = model.UserByEmail(db, cac, req.Email); err != nil {
		res.Result = 500
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if u.Salted != model.Stretch(req.Password, u.Salt) {
		res.Result = 401
		res.Error = "auth failed"
		c.JSON(http.StatusOK, res)
		return
	}
	u.LastLogin = ts.Unix()
	u.Updated = ts.Unix()
	if _, err := u.Update(db, cac); err != nil {
		res.Result = 500
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Result = 200
	res.User = u
	if err := logg.Log("user.login", res.User, ts); err != nil {
		log.Printf("post to fluentd failed %s", err)
	}

	c.JSON(http.StatusOK, res)
	return
}
