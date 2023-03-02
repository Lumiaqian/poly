package app

import (
	"changeme/internal/focus"
	"changeme/internal/liveroom"
	"changeme/internal/platform"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

type Server struct {
	log          *wails.CustomLogger
	focusService focus.FcousService
	huya         platform.HuYa
	douyu        platform.DouYu
	bilibili     platform.Bilibili
}

func NewServer() *Server {
	return &Server{
		log:          logger.NewCustomLogger("Server"),
		focusService: focus.NewFcousService(),
		huya:         platform.NewHuYa(),
		douyu:        platform.NewDoYu(),
		bilibili:     platform.NewBilibili(),
	}
}

// 获取直播流
func (s *Server) GetLiveRoom(platformName, roomId string) liveroom.LiveRoom {
	s.log.InfoFields("GetLiveRoom", logger.Fields{"platformName": platformName})
	room := liveroom.LiveRoom{}
	switch platformName {
	case platform.Huya:
		room, err := s.huya.GetLiveUrl(roomId)
		if err != nil {
			s.log.InfoFields("huya.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		s.log.InfoFields("huya.GetLiveUrl", logger.Fields{"room": room})
		return *room
	case platform.Bili:
		room, err := s.bilibili.GetLiveUrl(roomId)
		if err != nil {
			s.log.InfoFields("bilibili.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		s.log.InfoFields("bilibili.GetLiveUrl", logger.Fields{"room": room})
		return *room
	case platform.Douyu:
		room, err := s.douyu.GetLiveUrl(roomId)
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
func (a *Server) GetLiveRoomInfo(platformName, roomId string) liveroom.LiveRoomInfo {
	roomInfo := liveroom.LiveRoomInfo{}
	a.log.InfoFields("GetLiveRoomInfo", logger.Fields{"platformName": platformName})
	switch platformName {
	case platform.Huya:
		roomInfo, err := a.huya.GetRoomInfo(roomId)
		if err != nil {
			a.log.ErrorFields("huya.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	case platform.Bili:
		roomInfo, err := a.bilibili.GetRoomInfo(roomId)
		if err != nil {
			a.log.ErrorFields("bilibili.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	case platform.Douyu:
		roomInfo, err := a.douyu.GetRoomInfo(roomId)
		if err != nil {
			a.log.ErrorFields("douyu.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	}
	return roomInfo
}

// 获取关注列表
func (a *Server) GetFocus() []liveroom.LiveRoomInfo {
	return a.focusService.GetFcousRoomInfo()
}
