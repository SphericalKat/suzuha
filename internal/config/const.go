package config

import (
	"net/url"
	"regexp"
	"time"

	"github.com/allegro/bigcache"
)

var Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(48 * time.Hour))
//var ShortTimeCache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Hour))

var MalUrl, _ = url.Parse("https://myanimelist.net")

var InfoLinkRe = regexp.MustCompile(`/(\w+)/(?:\w+/)?(\d+)/.*`)

var JSONContentType = []byte("application/json; charset=utf-8")
