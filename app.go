package main

import (
	"changeme/focus"
	"changeme/liveroom"
	"changeme/platform"
	"context"
	"fmt"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	log          *wails.CustomLogger
	focusService focus.FcousService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	a.log = logger.NewCustomLogger("APP")
	a.focusService = focus.NewFcousService()
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 获取直播流
func (a *App) GetLiveRoom(platformName, roomId string) liveroom.LiveRoom {
	a.log.InfoFields("GetLiveRoom", logger.Fields{"platformName": platformName})
	room := liveroom.LiveRoom{}
	switch platformName {
	case platform.Huya:
		huya := platform.NewHuYa()
		room, err := huya.GetLiveUrl(a.ctx, roomId)
		if err != nil {
			a.log.InfoFields("huya.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		a.log.InfoFields("huya.GetLiveUrl", logger.Fields{"room": room})
		return *room
	case platform.Bili:
		bilibili := platform.NewBilibili()
		room, err := bilibili.GetLiveUrl(a.ctx, roomId)
		if err != nil {
			a.log.InfoFields("bilibili.GetLiveUrl", logger.Fields{"error": err})
			return liveroom.LiveRoom{}
		}
		a.log.InfoFields("bilibili.GetLiveUrl", logger.Fields{"room": room})
		return *room
	}
	return room
}

// Dplay画质对象转ArtPlay画质对象
func (a *App) ChangeQualityFromD2Art(qualities []liveroom.Quality) []liveroom.ArtQuality {
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
func (a *App) GetLiveRoomInfo(platformName, roomId string) liveroom.LiveRoomInfo {
	roomInfo := liveroom.LiveRoomInfo{}
	a.log.InfoFields("GetLiveRoomInfo", logger.Fields{"platformName": platformName})
	switch platformName {
	case platform.Huya:
		huya := platform.NewHuYa()
		roomInfo, err := huya.GetRoomInfo(roomId)
		if err != nil {
			a.log.ErrorFields("huya.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	case platform.Bili:
		bilibili := platform.NewBilibili()
		roomInfo, err := bilibili.GetRoomInfo(roomId)
		if err != nil {
			a.log.ErrorFields("bilibili.GetRoomInfo Err", logger.Fields{"error": err})
			return roomInfo
		}
		return roomInfo
	}
	return roomInfo
}

// 获取关注列表
func (a *App) GetFocus() []liveroom.LiveRoomInfo {
	return a.focusService.GetFcousRoomInfo()
}

// SelectFile 选择配置文件
func (a *App) SelectFile() (string, error) {
	title := "选择文件"
	filetype := "*.yml;*.yaml"
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文本数据",
				Pattern:     filetype,
			},
		},
	})
	if err != nil {
		return "", err
	}
	a.log.InfoFields("SelectFile path", logger.Fields{"path": selection})
	return selection, nil
}

// 加载关注列表
func (a *App) LoadFocus() []liveroom.LiveRoomInfo {
	path, err := a.SelectFile()
	if err != nil {
		a.log.ErrorFields("LoadFocus SelectFile Err", logger.Fields{"err": err})
		return nil
	}
	err = a.focusService.InitFcous(path)
	//a.log.InfoFields("")
	if err != nil {
		a.MessageDialog("加载关注列表", "失败！")
		return nil
	}
	a.MessageDialog("加载关注列表", "成功！")
	list := a.GetFocus()
	a.log.InfoFields("加载关注列表 成功！", logger.Fields{"focusList": list})
	return list
}

// 消息弹窗
func (a *App) MessageDialog(title, message string) {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:         title,
		Message:       message,
		Buttons:       []string{"Yes"},
		DefaultButton: "Yes",
	})
	if err != nil {
		a.log.ErrorFields("MessageDialog Err", logger.Fields{"err": err})
	}

}
