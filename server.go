package main

import (
	"github.com/deletescape/suzuha/api/anime"
	"github.com/deletescape/suzuha/api/search"
	"github.com/deletescape/suzuha/api/season"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"os"
)

const indexString = "Suzuha v0.1\ngithub.com/deletescape/suzuha"

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
		ctx.SetStatusCode(http.StatusInternalServerError)
		log.Println("PANIC:", i)
	}

	base := os.Getenv("BASE_PATH")

	mux.GET(base+"/", Index)
	mux.GET(base+"/alive", Alive)
	mux.GET(base+"/anime/:id", anime.Index)
	mux.GET(base+"/search/autocomplete", search.Autocomplete)
	mux.GET(base+"/search/anime", search.Anime)
	mux.GET(base+"/season/:year?/:season?", season.Season)

	log.Println("Starting suzuha")
	log.Fatal(fasthttp.ListenAndServe(":8081", mux.Handler))
}
