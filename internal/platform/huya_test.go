package platform

import (
	"testing"
)

var huya = NewHuYa()

func TestGetStreamInfo(t *testing.T) {
	huya.GetStreamInfo("222523")
}

func TestGetLiveUrl(t *testing.T) {
	liveroom, err := huya.GetLiveUrl("222523")
	if err != nil {
		t.Fail()
	}
	t.Logf("liveroom:%+v", liveroom)
}

func TestGetRoomInfo(t *testing.T) {
	liveRoomInfo, err := huya.GetRoomInfo("333003")
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
	list, err := huya.GetRecommend(2, 12)
	if err != nil {
		t.Error(err)
	}
	t.Logf("list:%+v", list)
}
