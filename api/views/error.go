package views

import (
	"github.com/deletescape/suzuha/internal/config"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"log"
	"net/http"
)

type ErrView struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func Wrap(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetContentTypeBytes(config.JSONContentType)
	errView := ErrView{
		Message: err.Error(),
		Status:  http.StatusNotFound,
	}
	log.Println(errView)

	json, _ := jettison.Marshal(errView)
	_, _ = ctx.Write(json)
}
