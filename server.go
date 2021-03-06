package main

import (
	"github.com/deletescape/suzuha/api/anime"
	"github.com/deletescape/suzuha/api/person"
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

	// General
	mux.GET(base+"/", Index)
	mux.GET(base+"/alive", Alive)

	// Anime
	mux.GET(base+"/anime/:id", anime.Index)

	// Search
	mux.GET(base+"/search/autocomplete", search.Autocomplete)
	mux.GET(base+"/search/anime", search.Anime)

	// Season
	mux.GET(base+"/season/:year?/:season?", season.Season)

	// Person
	mux.GET(base+"/person/:id", person.Index)

	log.Println("Starting suzuha")
	log.Fatal(fasthttp.ListenAndServe(":8081", mux.Handler))
}
