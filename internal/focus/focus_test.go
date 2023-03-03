package focus

import "testing"

var focusService = NewFocusService()

func TestGetFocusRoomInfo(t *testing.T) {
	focusService.InitFocus("/Users/lumiaqian/Desktop/focus.yml")
	list := focusService.GetFocusRoomInfo()
	if list == nil {
		t.Fail()
	}
	t.Logf("roomList: %+v", list)
}
