package search

import (
	"fmt"
	"github.com/deletescape/toraberu/internal/config"
	"github.com/deletescape/toraberu/pkg/scraper/search"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"net/http"
	"strconv"
)

func Anime(ctx *fasthttp.RequestCtx) {
	if !ctx.QueryArgs().Has("q") {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	query := string(ctx.QueryArgs().Peek("q"))
	if len(query) < 3 {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	page := 1
	var err error
	if ctx.QueryArgs().Has("page") {
		page, err = strconv.Atoi(string(ctx.QueryArgs().Peek("page")))
		if err != nil {
			ctx.SetStatusCode(http.StatusBadRequest)
			return
		}
	}
	cacheKey := fmt.Sprintf("search:anime:%s:%d", query, page)

	json, err := config.Cache.Get(cacheKey)
	if err != nil {
		animes, err := search.ScrapeAnimeSearch(query, page)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		json, err = jettison.Marshal(animes)
		if err != nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}
		go config.Cache.Set(cacheKey, json)
	}
	ctx.Write(json)
	ctx.SetContentTypeBytes(config.JSONContentType)
}
