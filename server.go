package main

import (
	"github.com/fasthttp/router"
	"github.com/deletescape/toraberu/api/anime"
	"github.com/deletescape/toraberu/api/search"
	"github.com/deletescape/toraberu/internal/config"
	"github.com/deletescape/toraberu/pkg/entities"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"log"
)

var indexInfo []byte

func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentTypeBytes(config.JSONContentType)
	_, _ = ctx.Write(indexInfo)
}

func main() {
	indexInfo, _ = jettison.Marshal(entities.IndexInfo{
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

	mux := router.New()

	mux.GET("/", Index)
	mux.GET("/anime/:id", anime.Index)
	mux.GET("/search/autocomplete", search.Autocomplete)
	mux.GET("/search/anime", search.Anime)

	log.Println("Starting toraberu")
	log.Fatal(fasthttp.ListenAndServe(":8081", mux.Handler))
}
