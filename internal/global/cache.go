package global

import (
	"changeme/internal/liveroom"
	"changeme/pkg/cache"
	"strings"
)

var Cache = cache.Init()

var FocusMap map[string]liveroom.LiveRoomInfo = make(map[string]liveroom.LiveRoomInfo)

func FormatKey(biz string, keys ...string) string {
	return biz + ":" + strings.Join(keys, ":")
}
