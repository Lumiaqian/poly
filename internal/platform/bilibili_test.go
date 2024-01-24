package platform

import (
	"context"
	"testing"
)

var bilibili = NewBilibili()

func TestGetBilibiliLiveUrl(t *testing.T) {
	bilibili.GetLiveUrl(context.Background(), "7777")
}

func TestGetBilibiliLiveRoomInfo(t *testing.T) {
	info, err := bilibili.GetRoomInfo(context.Background(), "545068")
	if err != nil {
		t.Fail()
	}
	t.Logf("info:%+v", info)
}

func TestBiliGetRecommend(t *testing.T) {
	list, err := bilibili.GetRecommend(context.Background(), 1, 12)
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}
