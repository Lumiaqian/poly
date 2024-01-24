package app

import (
	"changeme/internal/focus"
	"changeme/internal/liveroom"
	"changeme/internal/platform"
	"context"
	"math"
	"sync"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

type Server struct {
	log          *wails.CustomLogger
	focusService focus.FocusService
	huya         platform.HuYa
	douyu        platform.DouYu
	bilibili     platform.Bilibili
	ctx          context.Context
}

func NewServer() *Server {
	return &Server{
		log:          logger.NewCustomLogger("Server"),
		focusService: focus.NewFocusService(),
		huya:         platform.NewHuYa(),
		douyu:        platform.NewDoYu(),
		bilibili:     platform.NewBilibili(),
		ctx:          nil,
	}
}

func (s *Server) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// 获取直播流
func (s *Server) GetLiveRoom(platformName, roomId string) liveroom.LiveRoom {
	s.log.InfoFields("GetLiveRoom", logger.Fields{"platformName": platformName, "roomId": roomId})
	room := liveroom.LiveRoom{}
	switch platformName {
	case platform.Huya:
		room, err := s.huya.GetLiveUrl(s.ctx, roomId)
		if err != nil {
			s.log.InfoFields("huya.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		s.log.InfoFields("huya.GetLiveUrl", logger.Fields{"room": room})
		return *room
	case platform.Bili:
		room, err := s.bilibili.GetLiveUrl(s.ctx, roomId)
		if err != nil {
			s.log.InfoFields("bilibili.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		s.log.InfoFields("bilibili.GetLiveUrl", logger.Fields{"room": room})
		return *room
	case platform.Douyu:
		room, err := s.douyu.GetLiveUrl(s.ctx, roomId)
		if err != nil {
			s.log.InfoFields("douyu.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		s.log.InfoFields("douyu.GetLiveUrl", logger.Fields{"room": room})
		return *room
	}
	return room
}

// Dplay画质对象转ArtPlay画质对象
func (a *Server) ChangeQualityFromD2Art(qualities []liveroom.Quality) []liveroom.ArtQuality {
	artQualities := make([]liveroom.ArtQuality, len(qualities))
	for i, quality := range qualities {
		if quality.Type == "" {
			continue
		}
		artQuality := liveroom.ArtQuality{
			Default: false,
			Html:    quality.Name,
			Url:     quality.Url,
		}
		if i == 0 {
			artQuality.Default = true
		}
		artQualities = append(artQualities, artQuality)
	}
	return artQualities
}

// 获取直播间详情
func (s *Server) GetLiveRoomInfo(platformName, roomId string) liveroom.LiveRoomInfo {
	roomInfo := liveroom.LiveRoomInfo{}
	s.log.InfoFields("GetLiveRoomInfo", logger.Fields{"platformName": platformName})
	switch platformName {
	case platform.Huya:
		roomInfo, err := s.huya.GetRoomInfo(s.ctx, roomId)
		if err != nil {
			s.log.ErrorFields("huya.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	case platform.Bili:
		roomInfo, err := s.bilibili.GetRoomInfo(s.ctx, roomId)
		if err != nil {
			s.log.ErrorFields("bilibili.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	case platform.Douyu:
		roomInfo, err := s.douyu.GetRoomInfo(s.ctx, roomId)
		if err != nil {
			s.log.ErrorFields("douyu.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	}
	return roomInfo
}

// 获取关注列表
func (s *Server) GetFocus() []liveroom.LiveRoomInfo {
	return s.focusService.GetFocusRoomInfo(s.ctx)
}

// 获取All推荐列表
func (s *Server) GetRecommend(page, pageSize int) []liveroom.LiveRoomInfo {
	roomInfos := make([]liveroom.LiveRoomInfo, 0, pageSize)
	wg := sync.WaitGroup{}
	ch := make(chan []liveroom.LiveRoomInfo)
	wg.Add(3)

	go func(page, pageSize int) {
		defer wg.Done()
		list, err := s.huya.GetRecommend(s.ctx, page, pageSize)
		if err != nil {
			return
		}
		ch <- list
	}(page, pageSize)

	go func(page, pageSize int) {
		defer wg.Done()
		list, err := s.douyu.GetRecommend(s.ctx, page, pageSize)
		if err != nil {
			return
		}
		ch <- list
	}(page, pageSize)

	go func(page, pageSize int) {
		defer wg.Done()
		list, err := s.bilibili.GetRecommend(s.ctx, page, pageSize)
		if err != nil {
			return
		}
		ch <- list
	}(page, pageSize)

	go func() {
		wg.Wait()
		close(ch)
	}()
	for list := range ch {
		roomInfos = append(roomInfos, list[0:int(math.Ceil(float64(pageSize)/3.0))]...)
	}
	roomInfos = roomInfos[0:pageSize]
	//a.log.InfoFields("GetRecommend roomInfos", logger.Fields{"roomInfos": roomInfos})
	return roomInfos
}
