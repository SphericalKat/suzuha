package main

import (
	"github.com/deletescape/suzuha/api/anime"
	"github.com/deletescape/suzuha/api/search"
	"github.com/deletescape/suzuha/api/season"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
)

const indexString  = "Suzuha v0.1\ngithub.com/deletescape/suzuha"
var indexBytes = []byte(indexString)
var ok = []byte{'O', 'K'}

func Index(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.Write(indexBytes)
}

func Alive(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.Write(ok)
}

func main() {
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
