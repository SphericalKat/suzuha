package search

import (
	"fmt"
	"github.com/deletescape/suzuha/internal/config"
	"github.com/valyala/fasthttp"
)

func Autocomplete(ctx *fasthttp.RequestCtx) {
	if !ctx.QueryArgs().Has("q") {
		ctx.SetStatusCode(400)
		return
	}
	query := string(ctx.QueryArgs().Peek("q"))
	if query == "" {
		ctx.SetStatusCode(400)
		return
	}
	searchType := "all"
	if ctx.QueryArgs().Has("type") {
		searchType = string(ctx.QueryArgs().Peek("type"))
	}
	cacheKey := fmt.Sprintf("search:autocomplete:%s:%s", searchType, query)
	var data []byte
	var err error
	data, err = config.Cache.Get(cacheKey)
	if err != nil {
		status := 0
		status, data, err = fasthttp.Get([]byte{}, fmt.Sprintf("https://myanimelist.net/search/prefix.json?type=%s&keyword=%s&v=1", searchType, query))
		if err == nil && status >= 200 && status < 400 {
			_ = config.Cache.Set(cacheKey, data)
		} else {
			if status != 0 {
				ctx.SetStatusCode(status)
			} else {
				ctx.SetStatusCode(500)
			}
			return
		}
	}
	ctx.Write(data)
	ctx.SetContentTypeBytes(config.JSONContentType)
}
