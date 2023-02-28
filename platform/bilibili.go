package platform

import (
	"changeme/liveroom"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

type Bilibili struct {
	httpClient http.Client
	log        *wails.CustomLogger
}

type RoomInit struct {
	RoomId     string
	ShortId    string
	Uid        string
	LiveStatus int
}

const (
	BilibiliFD  = 80    //流畅
	BilibiliLD  = 150   //高清
	BilibiliSD  = 250   //超清
	BilibiliHD  = 400   //蓝光
	BilibiliOD  = 10000 //原画
	Bilibili4K  = 20000 //4K
	BilibiliHDR = 30000 //杜比

	Platform = "bilibili"
)

func NewBilibili() Bilibili {
	return Bilibili{
		httpClient: http.Client{},
		log:        logger.NewCustomLogger("bilibili"),
	}
}

// 获取哔哩哔哩直播的真实流媒体地址，默认获取直播间提供的最高画质
func (b *Bilibili) GetLiveUrl(ctx context.Context, roomId string) (*liveroom.LiveRoom, error) {
	//先获取直播状态和真实房间号
	roomInit, err := b.getRealRid(roomId)
	if err != nil {
		return nil, errors.Wrap(err, "getRealRid err")
	}
	b.log.InfoFields("roomInit结果：", logger.Fields{"roomInit": roomInit})
	if roomInit.LiveStatus == 0 {
		return nil, errors.New("未开播或直播间不存在")
	}
	realUrl, err := b.GetUrl(roomInit.RoomId, BilibiliOD)
	if err != nil {
		return nil, errors.Wrap(err, "GetUrl err")
	}
	b.log.InfoFields("realUrl结果：", logger.Fields{"realUrl": realUrl})
	room := new(liveroom.LiveRoom)
	room.LiveUrl = realUrl
	room.Platform = Platform
	room.PlatformName = liveroom.GetPlatform(room.Platform)
	room.RoomId = roomId
	return room, nil
}

// 获取直播状态和真实房间号
func (b *Bilibili) getRealRid(roomId string) (RoomInit, error) {
	roomInit := RoomInit{}
	url := "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + roomId
	result, err := b.HttpGet(url)
	if err != nil {
		return roomInit, err
	}
	if err != nil {
		return roomInit, err
	}
	b.log.InfoFields("getRealRid结果：", logger.Fields{"result": string(result)})
	parse := gjson.Parse(string(result))
	roomInit.LiveStatus = int(parse.Get("data.live_status").Int())
	roomInit.RoomId = parse.Get("data.room_id").String()
	roomInit.ShortId = parse.Get("data.short_id").String()
	roomInit.Uid = parse.Get("data.uid").String()
	return roomInit, nil
}

// 获取哔哩哔哩直播的真实流媒体地址，默认获取直播间提供的最高画质
func (b *Bilibili) GetUrl(roomId string, qn int) (string, error) {
	baseUrl := "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo?" +
		"room_id=" + roomId +
		"&protocol=0,1" +
		"&format=0,1,2" +
		"&codec=0,1" +
		"&platform=h5" +
		"&ptype=8"
	url := baseUrl + "&qn=" + strconv.Itoa(qn)
	result, err := b.HttpGet(url)
	if err != nil {
		return "", err
	}
	b.log.InfoFields("getRoomPlayInfo结果：", logger.Fields{"result": string(result)})
	parse := gjson.Parse(string(result))
	streamInfo := parse.Get("data.playurl_info.playurl.stream").Array()
	qnMax := qn
	for _, data := range streamInfo {
		acceptQn := data.Get("format").Array()[0].Get("codec").Array()[0].Get("accept_qn").Array()
		for _, qn := range acceptQn {
			if qn.Int() > int64(qnMax) {
				qnMax = int(qn.Int())
			}
		}
	}
	b.log.InfoFields("qnMax结果：", logger.Fields{"qnMax": qnMax})
	if qnMax != qn {
		result, err = b.HttpGet(baseUrl + "&qn=" + strconv.Itoa(qnMax))
		if err != nil {
			return "", err
		}
		b.log.InfoFields("再次getRoomPlayInfo结果：", logger.Fields{"result": string(result)})
		parse = gjson.Parse(string(result))
		streamInfo = parse.Get("data.playurl_info.playurl.stream").Array()
	}
	b.log.InfoFields("streamInfo结果：", logger.Fields{"streamInfo": streamInfo})
	stramUrls := make(map[string]string, 0)
	urls := make([]string, 0)
	for _, data := range streamInfo {
		format := data.Get("format").Array()
		formatName := format[0].Get("format_name").String()
		if formatName == "ts" {
			base := format[len(format)-1].Get("codec").Array()[0].Get("base_url").String()
			urlInfo := format[len(format)-1].Get("codec").Array()[0].Get("url_info").Array()
			for i, info := range urlInfo {
				host := info.Get("host").String()
				extra := info.Get("extra").String()
				stramUrls["线路"+strconv.Itoa(i+1)] = host + base + extra
				urls = append(urls, host+base+extra)
			}
		}
	}
	b.log.InfoFields("stramUrls结果：", logger.Fields{"stramUrls": stramUrls, "urls": urls})
	return urls[0], nil
}

// Get请求
func (b *Bilibili) HttpGet(url string) (string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36")
	response, err := b.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	b.log.InfoFields("HttpGet结果：", logger.Fields{"result": string(result)})
	return string(result), nil
}

func (b *Bilibili) GetRoomInfo(roomId string) (liveroom.LiveRoomInfo, error) {
	roomInfo := liveroom.LiveRoomInfo{}
	url := "https://api.live.bilibili.com/xlive/web-room/v1/index/" +
		"getH5InfoByRoom?room_id=" + roomId
	result, err := b.HttpGet(url)
	if err != nil {
		return roomInfo, err
	}
	b.log.InfoFields("getH5InfoByRoom结果：", logger.Fields{"result": string(result)})
	parse := gjson.Parse(string(result))
	info := parse.Get("data.room_info")
	anchorInfo := parse.Get("data.anchor_info.base_info")

	roomInfo.Platform = Platform
	roomInfo.PlatformName = liveroom.GetPlatform(roomInfo.Platform)
	roomInfo.RoomId = roomId
	roomInfo.RoomName = info.Get("title").String()
	roomInfo.OnLineCount = int(info.Get("online").Int())
	if int(info.Get("live_status").Int()) == 1 {
		roomInfo.LiveStatus = 2
	} else {
		roomInfo.LiveStatus = 0
	}
	roomInfo.Screenshot = info.Get("cover").String()
	if roomInfo.Screenshot != "" && !strings.Contains(roomInfo.Screenshot, "https") {
		roomInfo.Screenshot = strings.ReplaceAll(roomInfo.Screenshot, "http", "https")
	}
	roomInfo.Anchor = anchorInfo.Get("uname").String()
	roomInfo.Avatar = anchorInfo.Get("face").String()
	if roomInfo.Avatar != "" && !strings.Contains(roomInfo.Avatar, "https") {
		roomInfo.Avatar = strings.ReplaceAll(roomInfo.Avatar, "http", "https")
	}
	roomInfo.GameFullName = anchorInfo.Get("area_name").String()
	return roomInfo, nil
}
