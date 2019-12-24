package config

import (
	"github.com/allegro/bigcache"
	"time"
)

var Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(12 * time.Hour))
var JsonContentType = []byte("application/json")