package focus

import (
	"context"
	"testing"
)

var focusService = NewFocusService()

func TestGetFocusRoomInfo(t *testing.T) {
	focusService.InitFocus(context.TODO(), "/Users/lumiaqian/Desktop/focus.yml")
	list := focusService.GetFocusRoomInfo(context.TODO())
	if list == nil {
		t.Fail()
	}
	t.Logf("roomList: %+v", list)
}
