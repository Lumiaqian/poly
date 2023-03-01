package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

func Init() *cache.Cache {
	return cache.New(cache.NoExpiration, 1*time.Minute)
}
