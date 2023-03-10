package focus

import (
	"changeme/internal/global"
	"changeme/internal/liveroom"
	"changeme/internal/platform"
	"changeme/pkg/file"
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
	FocusName = "Focus"
)

type FocusList struct {
	Focus []Item `yaml:"focus"`
}

type Item struct {
	Platform string `yaml:"platform"` //直播平台
	RoomId   string `yaml:"roomId"`   //房间ID
	Anchor   string `yaml:"anchor"`   //主播
}

type FocusService struct {
	huya      platform.HuYa
	bilibili  platform.Bilibili
	douyu     platform.DouYu
	log       *wails.CustomLogger
	focusList FocusList
	roomList  []liveroom.LiveRoomInfo
	cache     cache.Cache
	wg        sync.WaitGroup
}

func NewFocusService() FocusService {
	log := logger.NewCustomLogger("focus")
	return FocusService{
		huya:      platform.NewHuYa(),
		bilibili:  platform.NewBilibili(),
		douyu:     platform.NewDoYu(),
		log:       log,
		focusList: FocusList{},
		roomList:  []liveroom.LiveRoomInfo{},
		cache:     *global.Cache,
		wg:        sync.WaitGroup{},
	}
}

func (f *FocusService) InitFocus(path string) error {
	if len(f.roomList) > 0 {
		f.roomList = f.roomList[0:0]
	}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &f.focusList)
	if err != nil {
		return err
	}
	f.getFocusRoomInfo()
	return nil
}

func (f *FocusService) GetFocusRoomInfo() []liveroom.LiveRoomInfo {
	return f.getFocusRoomInfo()
}

func (f *FocusService) getFocusRoomInfo() []liveroom.LiveRoomInfo {
	f.roomList = f.roomList[0:0]

	if len(global.FocusMap) > 0 {
		f.focusList.Focus = f.focusList.Focus[0:0]
		for _, val := range global.FocusMap {
			f.focusList.Focus = append(f.focusList.Focus, Item{
				Platform: val.Platform,
				RoomId:   val.RoomId,
				Anchor:   val.Anchor,
			})
		}
	}

	ch := make(chan *liveroom.LiveRoomInfo)
	f.wg.Add(len(f.focusList.Focus))
	for _, focus := range f.focusList.Focus {
		go f.getRoomInfo(focus, ch)
	}

	go func() {
		f.wg.Wait()
		close(ch)
	}()
	for info := range ch {
		f.roomList = append(f.roomList, *info)
		global.FocusMap[global.FormatKey(liveroom.FocusKey, info.Platform, info.RoomId)] = *info
	}
	sort.Sort(liveroom.LiveRoomInfoArray(f.roomList))
	return f.roomList
}

func (f *FocusService) getRoomInfo(focus Item, ch chan *liveroom.LiveRoomInfo) {
	defer f.wg.Done()
	if roomInfo, ok := global.Cache.Get(global.FormatKey(FocusName, focus.Platform, focus.RoomId)); ok {
		ch <- roomInfo.(*liveroom.LiveRoomInfo)
		return
	}
	switch focus.Platform {
	case platform.Huya:
		roomInfo, err := f.huya.GetRoomInfo(focus.RoomId)
		if err != nil {
			f.log.ErrorFields("GetRoomInfo Huya Err", logger.Fields{"err": err})
			return
		}
		global.Cache.Set(global.FormatKey(FocusName, focus.Platform, focus.RoomId), &roomInfo, 3*time.Minute)
		ch <- &roomInfo
	case platform.Bili:
		roomInfo, err := f.bilibili.GetRoomInfo(focus.RoomId)
		if err != nil {
			f.log.ErrorFields("GetRoomInfo Bilibili Err", logger.Fields{"err": err})
			return
		}
		global.Cache.Set(global.FormatKey(FocusName, focus.Platform, focus.RoomId), &roomInfo, 3*time.Minute)
		ch <- &roomInfo
	case platform.Douyu:
		roomInfo, err := f.douyu.GetRoomInfo(focus.RoomId)
		if err != nil {
			f.log.ErrorFields("GetRoomInfo douyu Err", logger.Fields{"err": err})
			return
		}
		global.Cache.Set(global.FormatKey(FocusName, focus.Platform, focus.RoomId), &roomInfo, 3*time.Minute)
		ch <- &roomInfo
	}
}

func (f *FocusService) Save(roomInfo liveroom.LiveRoomInfo) {
	if _, ok := global.FocusMap[global.FormatKey(liveroom.FocusKey, roomInfo.Platform, roomInfo.RoomId)]; ok {
		return
	}
	global.FocusMap[global.FormatKey(liveroom.FocusKey, roomInfo.Platform, roomInfo.RoomId)] = roomInfo
}

func (f *FocusService) Remove(roomInfo liveroom.LiveRoomInfo) {
	delete(global.FocusMap, global.FormatKey(liveroom.FocusKey, roomInfo.Platform, roomInfo.RoomId))
}

func (f *FocusService) SaveFocus() {
	f.focusList.Focus = f.focusList.Focus[0:0]
	for _, val := range global.FocusMap {
		f.focusList.Focus = append(f.focusList.Focus, Item{
			Platform: val.Platform,
			RoomId:   val.RoomId,
			Anchor:   val.Anchor,
		})
	}
	f.log.Info("保存关注的列表到文件")
	path := "./config/"
	fileName := "focus.yml"
	data, err := yaml.Marshal(f.focusList)
	if err != nil {
		f.log.ErrorFields("yml文件转换失败", logger.Fields{"err": err})
	}
	err = file.CreateFileWithDir(path, fileName, data)
	if err != nil {
		f.log.ErrorFields("yml文件写入失败", logger.Fields{"err": err})
	}
}
