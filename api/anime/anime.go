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
		id, _ := strconv.Atoi(idStr)
		json, _ = jettison.Marshal(anime.ScrapeAnime(id))
		config.Cache.Set(cacheKey, json)
	}

	ctx.SetContentTypeBytes(config.JsonContentType)
	ctx.Write(json)
}
