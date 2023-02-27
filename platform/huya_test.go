package platform

import (
	"context"
	"testing"
)

var huya = NewHuYa()

func TestGetStreamInfo(t *testing.T) {
	huya.GetStreamInfo("222523")
}

func TestGetLiveUrl(t *testing.T) {
	liveroom, err := huya.GetLiveUrl(context.Background(), "222523")
	if err != nil {
		t.Fail()
	}
	t.Logf("liveroom:%+v", liveroom)
}

func TestGetRoomInfo(t *testing.T) {
	liveRoomInfo, err := huya.GetRoomInfo("222523")
	if err != nil {
		t.Fail()
	}
	t.Logf("liveRoomInfo:%+v", liveRoomInfo)
}
