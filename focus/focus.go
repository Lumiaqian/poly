package focus

import (
	"changeme/global"
	"changeme/liveroom"
	"changeme/platform"
	"io/ioutil"
	"sort"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
	"gopkg.in/yaml.v3"
)

const (
	FcousName = "Fcous"
)

type FcousList struct {
	Fcous []Item `yaml:"fcous"`
}

type Item struct {
	Platform string `yaml:"platform"` //直播平台
	RoomId   string `yaml:"roomId"`   //房间ID
	Anchor   string `yaml:"anchor"`   //主播
}

type FcousService struct {
	huya      platform.HuYa
	bilibili  platform.Bilibili
	douyu     platform.DouYu
	log       *wails.CustomLogger
	fcousList FcousList
	roomList  []liveroom.LiveRoomInfo
	cache     cache.Cache
	wg        sync.WaitGroup
}

func NewFcousService() FcousService {
	log := logger.NewCustomLogger("fcous")
	return FcousService{
		huya:      platform.NewHuYa(),
		bilibili:  platform.NewBilibili(),
		douyu:     platform.NewDoYu(),
		log:       log,
		fcousList: FcousList{},
		roomList:  []liveroom.LiveRoomInfo{},
		cache:     *global.Cache,
		wg:        sync.WaitGroup{},
	}
}

func (f *FcousService) InitFcous(path string) error {
	if len(f.roomList) > 0 {
		f.roomList = f.roomList[0:0]
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &f.fcousList)
	if err != nil {
		return err
	}
	f.getFcousRoomInfo()
	return nil
}

func (f *FcousService) GetFcousRoomInfo() []liveroom.LiveRoomInfo {
	return f.getFcousRoomInfo()
}

func (f *FcousService) getFcousRoomInfo() []liveroom.LiveRoomInfo {
	f.roomList = f.roomList[0:0]

	ch := make(chan *liveroom.LiveRoomInfo)
	f.wg.Add(len(f.fcousList.Fcous))
	for _, fcous := range f.fcousList.Fcous {
		go f.getRoomInfo(fcous, ch)
	}

	go func() {
		f.wg.Wait()
		close(ch)
	}()
	for info := range ch {
		f.roomList = append(f.roomList, *info)
	}
	sort.Sort(liveroom.LiveRoomInfoArray(f.roomList))
	return f.roomList
}

func (f *FcousService) getRoomInfo(fcous Item, ch chan *liveroom.LiveRoomInfo) {
	defer f.wg.Done()
	if roomInfo, ok := global.Cache.Get(global.FormatKey(FcousName, fcous.Platform, fcous.RoomId)); ok {
		ch <- roomInfo.(*liveroom.LiveRoomInfo)
		return
	}
	switch fcous.Platform {
	case platform.Huya:
		roomInfo, err := f.huya.GetRoomInfo(fcous.RoomId)
		if err != nil {
			f.log.ErrorFields("GetRoomInfo Huya Err", logger.Fields{"err": err})
			return
		}
		global.Cache.Set(global.FormatKey(FcousName, fcous.Platform, fcous.RoomId), &roomInfo, 3*time.Minute)
		ch <- &roomInfo
	case platform.Bili:
		roomInfo, err := f.bilibili.GetRoomInfo(fcous.RoomId)
		if err != nil {
			f.log.ErrorFields("GetRoomInfo Bilibili Err", logger.Fields{"err": err})
			return
		}
		global.Cache.Set(global.FormatKey(FcousName, fcous.Platform, fcous.RoomId), &roomInfo, 3*time.Minute)
		ch <- &roomInfo
	case platform.Douyu:
		roomInfo, err := f.douyu.GetRoomInfo(fcous.RoomId)
		if err != nil {
			f.log.ErrorFields("GetRoomInfo douyu Err", logger.Fields{"err": err})
			return
		}
		global.Cache.Set(global.FormatKey(FcousName, fcous.Platform, fcous.RoomId), &roomInfo, 3*time.Minute)
		ch <- &roomInfo
	}
}
