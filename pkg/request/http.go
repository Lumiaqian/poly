package request

import (
	"net"
	"net/http"
	"time"

	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
)

var (
	dial = net.Dial
	tsp  = &http.Transport{Dial: dial}
)

// 发送HTTP请求
func HTTP() *dataflow.DataFlow {
	client := &http.Client{Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		DisableKeepAlives: true,
	},
		Timeout: 5 * time.Second}
	return gout.New(client).Debug(false)
}
