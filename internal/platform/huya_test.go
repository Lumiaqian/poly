package platform

import (
	"context"
	"testing"
)

var huya = NewHuYa()

func TestGetStreamInfo(t *testing.T) {
	huya.GetStreamInfo(context.Background(), "222523")
}

func TestGetLiveUrl(t *testing.T) {
	liveroom, err := huya.GetLiveUrl(context.Background(), "102411")
	if err != nil {
		t.Fail()
	}
	t.Logf("liveroom:%+v \n", liveroom)
	t.Logf("liveroom.LiveUrl:%+v\n", liveroom.LiveUrl)
}

func TestGetRoomInfo(t *testing.T) {
	liveRoomInfo, err := huya.GetRoomInfo(context.Background(), "333003")
	if err != nil {
		t.Fail()
	}
	t.Logf("liveRoomInfo:%+v", liveRoomInfo)
}

func TestGetSimgleArea(t *testing.T) {
	list, err := huya.GetSimgleArea("1", "网游竞技")
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}

func TestGetAllArea(t *testing.T) {
	list, err := huya.GetAllAreaInfo()
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}

func TestHuyaGetRecommend(t *testing.T) {
	list, err := huya.GetRecommend(context.Background(), 2, 12)
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}
