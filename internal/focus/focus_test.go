package focus

import "testing"

var fcousService = NewFcousService()

func TestGetFcousRoomInfo(t *testing.T) {
	fcousService.InitFcous("/Users/lumiaqian/Desktop/focus.yml")
	list := fcousService.GetFcousRoomInfo()
	if list == nil {
		t.Fail()
	}
	t.Logf("roomList: %+v", list)
}
