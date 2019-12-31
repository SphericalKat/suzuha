package main

import (
	"github.com/deletescape/suzuha/api/anime"
	"github.com/deletescape/suzuha/api/search"
	"github.com/deletescape/suzuha/api/season"
	"github.com/deletescape/suzuha/internal/config"
	"github.com/deletescape/suzuha/pkg/entities"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"log"
)

var indexInfo []byte
var alive = []byte("OK")

func Index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentTypeBytes(config.JSONContentType)
	_, _ = ctx.Write(indexInfo)
}

func Alive(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.Write(alive)
}

func main() {
	indexInfo, _ = jettison.Marshal(entities.IndexInfo{
		Author:           "@deletescape",
		Telegram:         "t.me/noneyet",
		Version:          "0.0.1",
		SuzuhaGo:       "0.0.1",
		Website:          "suzuha.deletescape.ch",
		Docs:             "suzuha.deletescape.ch/docs",
		GitHub:           "https://github.com/deletescape/suzuha",
		ProductionApiUrl: "https://suzuha.deletescape.cloud/api",
		StatusUrl:        "https://status.deletescape.cloud/suzuha",
	})

	mux := router.New()
	mux.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		ctx.SetStatusCode(500)
		log.Println("PANIC:", i)
	}

	mux.GET("/", Index)
	mux.GET("/alive", Alive)
	mux.GET("/anime/:id", anime.Index)
	mux.GET("/search/autocomplete", search.Autocomplete)
	mux.GET("/search/anime", search.Anime)
	mux.GET("/season/:year?/:season?", season.Season)

	log.Println("Starting suzuha")
	log.Fatal(fasthttp.ListenAndServe(":8081", mux.Handler))
}
