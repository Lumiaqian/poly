package platform

import (
	"testing"
)

var bilibili = NewBilibili()

func TestGetBilibiliLiveUrl(t *testing.T) {
	bilibili.GetLiveUrl("7777")
}

func TestGetBilibiliLiveRoomInfo(t *testing.T) {
	info, err := bilibili.GetRoomInfo("545068")
	if err != nil {
		t.Fail()
	}
	t.Logf("info:%+v", info)
}

func TestBiliGetRecommend(t *testing.T) {
	list, err := bilibili.GetRecommend(1, 12)
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}
