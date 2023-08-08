package app

import (
	C "changeme/internal/constant"
	"changeme/internal/liveroom"

	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) initMenu() {
	a.Menu = menu.NewMenu()
	fileMenu := a.Menu.AddSubmenu("File")
	fileMenu.AddText("配置", keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {
		a.LoadFocus()
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("退出", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(a.ctx)
	})
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
	err = a.Server.focusService.InitFocus(path)
	//a.log.InfoFields("")
	if err != nil {
		a.MessageDialog("加载关注列表", "失败！")
		return nil
	}
	a.MessageDialog("加载关注列表", "成功！")
	list := a.Server.GetFocus()
	a.log.InfoFields("加载关注列表 成功！", logger.Fields{"focusList": list})
	return list
}

// 加载指定位置的文件
func (a *App) LoadLocalFocus() []liveroom.LiveRoomInfo {
	path := C.Path.Focus()
	err := a.Server.focusService.InitFocus(path)
	//a.log.InfoFields("")
	if err != nil {
		a.MessageDialog("加载关注列表失败", "请从菜单->配置中选择文件加载")
		return nil
	}
	a.MessageDialog("加载关注列表", "成功！")
	list := a.Server.GetFocus()
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

// 关注
func (a *App) SaveFocus(roomInfo liveroom.LiveRoomInfo) {
	a.Server.focusService.Save(roomInfo)
}

// 移出关注
func (a *App) RemoveFocus(roomInfo liveroom.LiveRoomInfo) {
	a.Server.focusService.Remove(roomInfo)
}
