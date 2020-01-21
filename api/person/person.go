package person

import (
	"fmt"
	"github.com/deletescape/suzuha/api/views"
	"github.com/deletescape/suzuha/internal/config"
	"github.com/deletescape/suzuha/pkg/person"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
	"strconv"
)

func Index(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	cacheKey := fmt.Sprintf("person:%s", idStr)

	json, err := config.Cache.Get(cacheKey)
	if err != nil {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			views.Wrap(ctx, err)
			return
		}
		p, err := person.ScrapePerson(id)
		if err != nil {
			views.Wrap(ctx, err)
			return
		}
		json, err = jettison.Marshal(p)
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
