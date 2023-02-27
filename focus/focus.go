package focus

import (
	"changeme/liveroom"
	"changeme/platform"
	"io/ioutil"
	"sort"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
	"gopkg.in/yaml.v3"
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
	log       *wails.CustomLogger
	fcousList FcousList
	roomList  []liveroom.LiveRoomInfo
}

func NewFcousService() FcousService {
	log := logger.NewCustomLogger("fcous")
	return FcousService{
		huya:      platform.NewHuYa(),
		log:       log,
		fcousList: FcousList{},
		roomList:  []liveroom.LiveRoomInfo{},
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
	for _, fcous := range f.fcousList.Fcous {
		switch fcous.Platform {
		case "huya":
			roomInfo, err := f.huya.GetRoomInfo(fcous.RoomId)
			if err != nil {
				f.log.ErrorFields("GetRoomInfo Err", logger.Fields{"err": err})
				continue
			}
			f.roomList = append(f.roomList, roomInfo)
		}
	}
	sort.Sort(liveroom.LiveRoomInfoArray(f.roomList))
	return nil
}

func (f *FcousService) GetFcousRoomInfo() []liveroom.LiveRoomInfo {
	return f.roomList
}
