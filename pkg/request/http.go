package request

import (
	"net"
	"net/http"

	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
)

var dial = net.Dial

// 发送HTTP请求
func HTTP() *dataflow.DataFlow {
	client := http.DefaultClient
	tsp := &http.Transport{Dial: dial}
	client.Transport = tsp
	return gout.New(client).Debug(false)
}
