package global

import (
	"changeme/pkg/cache"
	"strings"
)

var Cache = cache.Init()

func FormatKey(biz string, keys ...string) string {
	return biz + ":" + strings.Join(keys, ":")
}
