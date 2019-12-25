package config

import (
	"time"

	"github.com/allegro/bigcache"
)

var Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(48 * time.Hour))

//var ShortTimeCache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Hour))
var JSONContentType = []byte("application/json; charset=utf-8")
