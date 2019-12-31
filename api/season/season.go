package season

import (
	"fmt"
	"github.com/deletescape/toraberu/internal/config"
	"github.com/deletescape/toraberu/pkg/scraper/season"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"net/http"
)

func Season(ctx *fasthttp.RequestCtx) {
	var year string
	var seas string
	var cacheKey string

	yearI := ctx.UserValue("year")
	if yearI != nil {
		year = yearI.(string)
		seasI := ctx.UserValue("season")
		if seasI != nil {
			seas = seasI.(string)
		} else {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		cacheKey = fmt.Sprintf("season:%s:%s", year, seas)
	} else {
		cacheKey = "season:current"
	}

	json, err := config.Cache.Get(cacheKey)
	if err != nil {
		season, err := season.ScrapeAnime(year, seas)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		json, err = jettison.Marshal(season)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		go config.Cache.Set(cacheKey, json)
	}
	ctx.Write(json)
	ctx.SetContentTypeBytes(config.JSONContentType)
}