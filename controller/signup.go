package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"simple-web-golang/model"
	"unicode"
)

type SignUpRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignUpResponse struct {
	Result int        `json:"result"`
	Error  string     `json:"error"`
	User   model.User `json:"user"`
}

func (h *Controller) SignUp(c *gin.Context) {
	var req SignUpRequest
	var res SignUpResponse
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ts, _ := TsFrom(c)
	cac, _ := CacheFrom(c)
	db, _ := DBFrom(c)
	logg, _ := LoggerFrom(c)

	b, err := model.UserExist(db, cac, req.Email)
	if err != nil {
		res.Result = 500
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if b {
		res.Result = 202
		res.Error = "given email address is already used"
		c.JSON(http.StatusOK, res)
		return
	}

	u := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Created:   ts.Unix(),
		Updated:   ts.Unix(),
		LastLogin: ts.Unix(),
	}
	if _, err := u.Insert(db, req.Password); err != nil {
		res.Result = 500
		res.Error = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if res.User, err = model.UserByEmail(db, cac, u.Email); err != nil {
		res.Result = 201
		res.Error = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	if err := logg.Log("user.created", res.User, ts); err != nil {
		log.Printf("post to fluentd failed %s", err)
	}

	c.JSON(http.StatusOK, res)
	return
}

func onlyLetter(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) {
			return false
		}
	}
	return true
}
