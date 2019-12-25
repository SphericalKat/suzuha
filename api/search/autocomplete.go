package search

import (
	"fmt"
	"github.com/deletescape/toraberu/config"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
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
		resp, err := http.Get(fmt.Sprintf("https://myanimelist.net/search/prefix.json?type=%s&keyword=%s&v=1", searchType, query))
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			data, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				ctx.SetStatusCode(500)
				return
			}
			config.Cache.Set(cacheKey, data)
		} else {
			if resp != nil {
				ctx.SetStatusCode(resp.StatusCode)
			} else {
				ctx.SetStatusCode(500)
			}
			return
		}
	}
	ctx.Write(data)
	ctx.SetContentTypeBytes(config.JSONContentType)
}