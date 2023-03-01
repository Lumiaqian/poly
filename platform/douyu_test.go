package platform

import (
	"context"
	"testing"
)

var douyu = NewDoYu()

func TestGetDouyuLiveUrl(t *testing.T) {
	liveUrl, err := douyu.GetLiveUrl(context.TODO(), "687423")
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
