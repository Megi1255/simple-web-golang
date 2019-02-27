package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"simple-web-golang/config"
	"simple-web-golang/model"
	"simple-web-golang/util"
)

type SearchArtistRequest struct {
	Mode   string `json:"mode" binding:"eq=|eq=name|eq=area|eq=alias|eq=tag"`
	Value  string `json:"value"`
	Offset int64  `json:"offset"`
	Sort   string `json:"sort" binding:"eq=rating|eq=name|eq="`
	SortBy int    `json:"sort_by" binding:"gte=-1,lte=1"`
}

type SearchArtistResponse struct {
	ApiName string         `json:"api_name"`
	Code    int            `json:"code"`
	Error   string         `json:"error"`
	Total   int64          `json:"total"`
	Offset  int64          `json:"offset"`
	Artists []model.Artist `json:"artists"`
}

func (h *HandlerFunc) SearchArtist(c *gin.Context) {
	var req SearchArtistRequest
	var res SearchArtistResponse
	res.ApiName = "SearchArtist"
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//ts := util.TsFromContext(c)
	//cac, _ := util.CacheFrom(c)
	//db, _ := util.MongoFrom(c)
	//logg := log.FromContext(c)
	prof := util.ProfilerFromContext(c)

	sortmap := map[string]string{
		"rating": "rating.votescount",
		"name":   "sortname",
	}
	if req.Sort != "" {
		req.Sort = sortmap[req.Sort]
	}

	var err error
	if req.Mode == "name" {
		prof.Start("db-query.ArtistsByName")
		if res.Total, res.Artists, err = model.ArtistsByName(c, req.Value, 10, req.Sort, req.SortBy); err != nil {
			HandleError(c, res.ApiName, config.CODE_DB_ERROR, err.Error())
			return
		}
		prof.End("db-query.ArtistsByName")
		if res.Total == 0 {
			prof.Start("db-query.ArtistsByAlias")
			if res.Total, res.Artists, err = model.ArtistsByAlias(c, req.Value, req.Offset, 10, req.Sort, req.SortBy); err != nil {
				HandleError(c, res.ApiName, config.CODE_DB_ERROR, err.Error())
				return
			}
			prof.End("db-query.ArtistsByAlias")
		}
	} else if req.Mode == "tag" {
		prof.Start("db-query.ArtistsByTag")
		if res.Total, res.Artists, err = model.ArtistByTag(c, req.Value, req.Offset, 10, req.Sort, req.SortBy); err != nil {
			HandleError(c, res.ApiName, config.CODE_DB_ERROR, err.Error())
			return
		}
		prof.End("db-query.ArtistsByTag")
	} else if req.Mode == "area" {
		prof.Start("db-query.ArtistsByArea")
		if res.Total, res.Artists, err = model.ArtistsByArea(c, req.Value, req.Offset, 10, req.Sort, req.SortBy); err != nil {
			HandleError(c, res.ApiName, config.CODE_DB_ERROR, err.Error())
			return
		}
		prof.End("db-query.ArtistsByArea")
	} else {
		res.Error = "not supported mode"
	}

	res.Code = config.CODE_OK
	res.Offset = req.Offset

	c.JSON(http.StatusOK, res)
	return
}
