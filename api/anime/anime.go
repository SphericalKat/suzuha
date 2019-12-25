package anime

import (
	"fmt"
	"github.com/deletescape/toraberu/api/views"
	"github.com/deletescape/toraberu/internal/config"
	"github.com/deletescape/toraberu/pkg/scraper/anime"
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
			views.Wrap(ctx, err)
			return
		}
		a, err := anime.ScrapeAnime(id)
		if err != nil {
			views.Wrap(ctx, err)
			return
		}
		json, err = jettison.Marshal(a)
		if err != nil {
			views.Wrap(ctx, err)
			return
		}
		err = config.Cache.Set(cacheKey, json)
		if err != nil {
			views.Wrap(ctx, err)
			return
		}
	}

	ctx.SetContentTypeBytes(config.JSONContentType)
	_, _ = ctx.Write(json)
}
