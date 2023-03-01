package platform

import (
	"changeme/liveroom"
	"changeme/pkg/request"
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/guonaihong/gout"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

/*
一共找到3个可用的CDN域名：
http://hw-tct.douyucdn.cn
http://hdltc1.douyucdn.cn
http://hdltctwk.douyucdn2.cn
rate: 1流畅；2高清；3超清；4蓝光4M；0蓝光8M或10M
*/
const (
	DouYuRoomId = `\$ROOM.room_id =(.*?); `
	DouYuMatch  = `(vdwdae325w_64we[\s\S]*function ub98484234[\s\S]*?)function`
)

type DouYu struct {
	httpClient http.Client
	log        *wails.CustomLogger
}

func NewDoYu() DouYu {
	return DouYu{
		httpClient: http.Client{},
		log:        logger.NewCustomLogger(Douyu),
	}
}

func (d *DouYu) GetLiveUrl(ctx context.Context, roomId string) (*liveroom.LiveRoom, error) {
	did := "10000000000000000000000000001501"
	t10 := strconv.FormatInt(time.Now().Unix(), 10)
	html := ""
	if err := request.HTTP().GET(fmt.Sprintf("https://www.douyu.com/%s", roomId)).BindBody(&html).Do(); err != nil {
		return nil, err
	}
	result := regexp.MustCompile(DouYuMatch).FindString(html)
	jsUb9 := strings.TrimSuffix(result, "function")
	jsUb9 = regexp.MustCompile(`eval.*?;}`).ReplaceAllString(jsUb9, `strc;}`)
	vm := goja.New()

	if _, err := vm.RunString(jsUb9); err != nil {
		return nil, err
	}

	ub9, ok := goja.AssertFunction(vm.Get("ub98484234"))
	if !ok {
		return nil, fmt.Errorf("failed to assert function ub9")
	}
	res, err := ub9(goja.Undefined())
	if err != nil {
		return nil, err
	}
	value := regexp.MustCompile(`v=(\d+)`).FindAllStringSubmatch(res.String(), -1)[0][1]

	rb := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%s%s", roomId, did, t10, value))))
	funcSign := regexp.MustCompile(`return rt;}\);?`).ReplaceAllString(res.String(), `return rt;}`)
	funcSign = strings.ReplaceAll(funcSign, `(function (`, `function sign(`)
	funcSign = strings.ReplaceAll(funcSign, `CryptoJS.MD5(cb).toString()`, `"`+rb+`"`)

	if _, err = vm.RunString(funcSign); err != nil {
		return nil, err
	}
	sign, ok := goja.AssertFunction(vm.Get("sign"))
	if !ok {
		return nil, fmt.Errorf("failed to assert function sign")
	}
	param, err := sign(goja.Undefined(), vm.ToValue(roomId), vm.ToValue(did), vm.ToValue(t10))
	if err != nil {
		return nil, err
	}
	params := fmt.Sprintf("%s&cdn=ws-h5&rate=%d", param, 0)

	var resp struct {
		Error int    `json:"error"`
		Msg   string `json:"msg"`
		Data  struct {
			RoomID       int64  `json:"room_id"`
			IsMixed      bool   `json:"is_mixed"`
			MixedLive    string `json:"mixed_live"`
			MixedURL     string `json:"mixed_url"`
			RtmpCdn      string `json:"rtmp_cdn"`
			RtmpURL      string `json:"rtmp_url"`
			RtmpLive     string `json:"rtmp_live"`
			ClientIP     string `json:"client_ip"`
			InNA         int    `json:"inNA"`
			RateSwitch   int    `json:"rateSwitch"`
			Rate         int    `json:"rate"`
			CdnsWithName []*struct {
				Name   string `json:"name"`
				Cdn    string `json:"cdn"`
				IsH265 bool   `json:"isH265"`
			} `json:"cdnsWithName"`
			Multirates []*struct {
				Name    string `json:"name"`
				Rate    int    `json:"rate"`
				HighBit int    `json:"highBit"`
				Bit     int    `json:"bit"`
			} `json:"multirates"`
		}
	}

	err = request.HTTP().POST(fmt.Sprintf("https://www.douyu.com/lapi/live/getH5Play/%s?", roomId) + params).
		SetHeader(gout.H{
			"UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
			"referer":   "https://www.douyu.com/",
			"origin":    "https://www.douyu.com",
		}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	d.log.InfoFields("resp", logger.Fields{"resp": resp})
	if resp.Error != 0 {
		return nil, fmt.Errorf("failed to get live url: %s", resp.Msg)
	}
	return &liveroom.LiveRoom{
		Platform:     Douyu,
		PlatformName: liveroom.GetPlatform(Douyu),
		RoomId:       roomId,
		LiveUrl:      fmt.Sprintf("https://hw-tct.douyucdn.cn/live/%s?uuid=", strings.Split(resp.Data.RtmpLive, "?")[0]),
	}, nil
}

func (d *DouYu) GetRoomInfo(roomId string) (liveroom.LiveRoomInfo, error) {
	roomInfo := liveroom.LiveRoomInfo{}
	var info struct {
		Error int `json:"error"`
		Data  struct {
			RoomId     string `json:"room_id"`
			OwnerName  string `json:"owner_name"`
			RoomStatus string `json:"room_status"`
			RoomName   string `json:"room_name"`
			CateName   string `json:"cate_name"`
			Avatar     string `json:"avatar"`
			Online     int    `json:"online"`
			RoomThumb  string `json:"room_thumb"`
		} `json:"data"`
	}
	if err := request.HTTP().GET(fmt.Sprintf("https://open.douyucdn.cn/api/RoomApi/room/%s", roomId)).BindJSON(&info).Do(); err != nil {

		return roomInfo, err
	}
	if info.Error != 0 {

		return roomInfo, errors.New("request err")
	}

	return liveroom.LiveRoomInfo{
		Platform:     Douyu,
		PlatformName: liveroom.GetPlatform(Douyu),
		RoomId:       roomId,
		RoomName:     info.Data.RoomName,
		Anchor:       info.Data.OwnerName,
		Avatar:       info.Data.Avatar,
		OnLineCount:  info.Data.Online,
		Screenshot:   info.Data.RoomThumb,
		GameFullName: info.Data.CateName,
		LiveStatus:   lo.If(info.Data.RoomStatus == "1", 2).Else(0),
	}, nil
}
