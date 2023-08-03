package platform

import (
	"testing"
)

var douyu = NewDoYu()

func TestGetRealUrl(t *testing.T) {
	room, err := douyu.GetRealUrl("7084089", "hls")
	if err != nil {
		t.Fail()
	}
	t.Logf("room:%+v", room)
}

func TestGetDouyuLiveUrl(t *testing.T) {
	liveUrl, err := douyu.GetLiveUrl("7084089")
	if err != nil {
		t.Fail()
	}
	t.Logf("liveUrl:%+v", liveUrl)
}

func TestGetDouyuRoomInfo(t *testing.T) {
	roomInfo, err := douyu.GetRoomInfo("687423")
	if err != nil {
		t.Fail()
	}
	t.Logf("roomInfo:%+v", roomInfo)
}

func TestDouyuGetRecommend(t *testing.T) {
	list, err := douyu.GetRecommend(2, 12)
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}
