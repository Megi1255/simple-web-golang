package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"simple-web-golang/model"
	"simple-web-golang/util"
)

type SearchArtistRequest struct {
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Area   string `json:"area"`
	Tags   string `json:"tags"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
	SortBy string `json:"sort_by"`
}

type SearchArtistResponse struct {
	Code    int            `json:"code"`
	Error   string         `json:"error"`
	Total   int            `json:"total"`
	Offset  int            `json:"offset"`
	Artists []model.Artist `json:"artists"`
}

func (h *Controller) SearchArtist(c *gin.Context) {
	var req SearchArtistRequest
	var res SearchArtistResponse
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ts, _ := util.TsFrom(c)
	cac, _ := util.CacheFrom(c)
	db, _ := util.DBFrom(c)
	logg, _ := util.LoggerFrom(c)

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
