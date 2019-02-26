package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"simple-web-golang/config"
	"simple-web-golang/model"
)

type SearchArtistRequest struct {
	Mode   string `json:"mode" binding:"eq=|eq=name|eq=area|eq=alias|eq=tag"`
	Value  string `json:"value"`
	Offset int64  `json:"offset"`
	Sort   string `json:"sort" binding:"eq=rating|eq=name|eq="`
	SortBy int    `json:"sort_by" binding:"gte=-1,lte=1"`
}

type SearchArtistResponse struct {
	Code    int            `json:"code"`
	Error   string         `json:"error"`
	Total   int            `json:"total"`
	Offset  int64          `json:"offset"`
	Artists []model.Artist `json:"artists"`
}

func (h *Controller) SearchArtist(c *gin.Context) {
	var req SearchArtistRequest
	var res SearchArtistResponse
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//ts := util.TsFromContext(c)
	//cac, _ := util.CacheFrom(c)
	//db, _ := util.MongoFrom(c)
	//logg := log.FromContext(c)

	sortmap := map[string]string{
		"rating": "rating.count",
		"name":   "sortname",
	}
	if req.Sort != "" {
		req.Sort = sortmap[req.Sort]
	}

	var err error
	if req.Mode == "name" {
		if res.Artists, err = model.ArtistsByName(c, req.Value, 20, req.Sort, req.SortBy); err != nil {
			HandleError(c, config.CODE_DB_ERROR, err.Error())
			return
		}
	} else if req.Mode == "tag" {
		if res.Artists, err = model.ArtistByTag(c, req.Value, req.Offset, 20, req.Sort, req.SortBy); err != nil {
			HandleError(c, config.CODE_DB_ERROR, err.Error())
			return
		}
	} else if req.Mode == "area" {
		if res.Artists, err = model.ArtistsByArea(c, req.Value, req.Offset, 20, req.Sort, req.SortBy); err != nil {
			HandleError(c, config.CODE_DB_ERROR, err.Error())
			return
		}
	} else {
		res.Error = "not supported mode"
	}

	res.Code = config.CODE_OK
	res.Offset = req.Offset

	c.JSON(http.StatusOK, res)
	return
}
