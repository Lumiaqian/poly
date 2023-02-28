package platform

import (
	"context"
	"testing"
)

var bilibili = NewBilibili()

func TestGetBilibiliLiveUrl(t *testing.T) {
	bilibili.GetLiveUrl(context.TODO(), "7777")
}

func TestGetBilibiliLiveRoomInfo(t *testing.T) {
	info, err := bilibili.GetRoomInfo("545068")
	if err != nil {
		t.Fail()
	}
	t.Logf("info:%+v", info)
}
