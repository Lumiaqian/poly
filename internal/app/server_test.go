package app

import (
	"testing"
)

var server = NewServer()

func TestGetRecommend(t *testing.T) {
	list := server.GetRecommend(2, 14)
	t.Logf("list:%+v,len:%d", list, len(list))
}
