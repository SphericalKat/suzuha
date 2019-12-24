package anime

import (
	"fmt"
	"github.com/deletescape/toraberu/config"
	"github.com/deletescape/toraberu/scraper/anime"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"strconv"
)

func Index(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	cacheKey := fmt.Sprintf("anime:%s", idStr)

	json, err := config.Cache.Get(cacheKey)
	if err != nil {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.SetStatusCode(404)
			fmt.Println(err)
			return
		}
		anime, err := anime.ScrapeAnime(id)
		if err != nil {
			ctx.SetStatusCode(404)
			fmt.Println(err)
			return
		}
		json, err = jettison.Marshal(anime)
		if err != nil {
			ctx.SetStatusCode(404)
			fmt.Println(err)
			return
		}
		config.Cache.Set(cacheKey, json)
	}

	ctx.SetContentTypeBytes(config.JsonContentType)
	_, _ = ctx.Write(json)
}
