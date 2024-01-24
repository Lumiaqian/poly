package platform

import (
	"changeme/internal/liveroom"
	"context"
)

const (
	Bili  = "bilibili"
	Huya  = "huya"
	Douyu = "douyu"
)

type PlatformRepo interface {
	GetLiveUrl(ctx context.Context, roomId string) (*liveroom.LiveRoom, error)
	GetRoomInfo(ctx context.Context, roomId string) (liveroom.LiveRoomInfo, error)
	GetRecommend(page, pageSize int) ([]liveroom.LiveRoomInfo, error)
}

// type Operator interface {
// 	GetLiveUrl(ctx context.Context, roomId string) (*liveroom.LiveRoom, error)
// 	GetRoomInfo(roomId string) (liveroom.LiveRoomInfo, error)
// }

// var (
// 	once    sync.Once
// 	factory = make(map[string]Operator)
// )

// func GetFactory(platform string)
