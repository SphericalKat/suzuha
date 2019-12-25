package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/deletescape/toraberu/api/anime"
	"github.com/deletescape/toraberu/api/search"
	"github.com/deletescape/toraberu/config"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"log"
)

type IndexInfo struct {
	Author           string
	Telegram         string
	Version          string
	ToraberuGo       string
	Website          string
	Docs             string
	GitHub           string
	ProductionApiUrl string `json:"PRODUCTION_API_URL"`
	StatusUrl        string `json:"STATUS_URL"`
}

var indexInfo []byte

func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentTypeBytes(config.JsonContentType)
	_, _ = ctx.Write(indexInfo)
}

func main() {
	indexInfo, _ = jettison.Marshal(IndexInfo{
		Author:           "@deletescape",
		Telegram:         "t.me/noneyet",
		Version:          "0.0.1",
		ToraberuGo:       "0.0.1",
		Website:          "toraberu.deletescape.ch",
		Docs:             "toraberu.deletescape.ch/docs",
		GitHub:           "https://github.com/deletescape/toraberu",
		ProductionApiUrl: "https://toraberu.deletescape.cloud/api",
		StatusUrl:        "https://status.deletescape.cloud/toraberu",
	})

	router := fasthttprouter.New()

	router.GET("/", Index)
	router.GET("/anime/:id", anime.Index)
	router.GET("/search/autocomplete", search.Autocomplete)
	router.GET("/search/anime", search.Anime)

	log.Fatal(fasthttp.ListenAndServe(":8081", router.Handler))
}
