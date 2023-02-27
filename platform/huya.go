package platform

import (
	"changeme/codec"
	"changeme/liveroom"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

type HuYa struct {
	httpClient http.Client
	log        *wails.CustomLogger
}

const (
	Nick          = `\"nick\":\"(.*?)\",`
	RoomName      = `\"roomName\":\"(.*?)\",`
	Screenshot    = `\"screenshot\":\"(.*?)\",`
	Avatar        = `\"avatar180\":\"(.*?)\",`
	GameFullName  = `\"gameFullName\":\"(.*?)\",`
	ActivityCount = `\"activityCount\":(.*?),`
	Introduction  = `\"introduction\":\"(.*?)\",`
	LiveStatus    = `\"liveStatus-(.*?) on-match\"`
)

func NewHuYa() HuYa {
	x := logger.NewCustomLogger("huya")
	return HuYa{
		httpClient: http.Client{Timeout: time.Second * 5},
		log:        x,
	}
}

// 获取真实直播流
func (h *HuYa) GetLiveUrl(ctx context.Context, roomId string) (*liveroom.LiveRoom, error) {
	roomUrl := "https://m.huya.com/" + roomId
	request, err := http.NewRequest("GET", roomUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36")
	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile("<script> window.HNF_GLOBAL_INIT = (.*)</script>")
	submatch := reg.FindStringSubmatch(string(result))
	if submatch == nil || len(submatch) < 2 {
		return nil, errors.Wrap(err, "查询失败！")
	}
	room, err := extractInfo(submatch[1])
	if err != nil {
		return nil, err
	}
	streamInfoList, err := h.GetStreamInfo(roomId)
	if err != nil {
		return nil, err
	}
	for _, streamInfo := range streamInfoList {
		quality := liveroom.Quality{
			Name: streamInfo.DisplayName,
			Url:  "",
			Type: "",
		}
		if streamInfo.BitRate != 0 {
			quality.Url = room.LiveUrl + "&ratio=" + strconv.FormatInt(streamInfo.BitRate, 10)
		} else {
			quality.Url = room.LiveUrl
		}
		quality.Type = GetType(quality.Url)
		room.Quality = append(room.Quality, quality)
	}
	roomInfo, err := h.GetRoomInfo(roomId)
	if err != nil {
		h.log.ErrorFields("GetRoomInfo Error", logger.Fields{"err": err, "roomInfo": roomInfo})
		return room, nil
	}
	room.Anchor = roomInfo.Anchor
	room.Avatar = roomInfo.Avatar
	room.RoomId = roomId
	room.OnLineCount = roomInfo.OnLineCount
	room.Platform = roomInfo.Platform
	room.Screenshot = roomInfo.Screenshot
	room.GameFullName = roomInfo.GameFullName
	room.PlatformName = liveroom.GetPlatform(room.Platform)
	return room, nil
}

func extractInfo(content string) (*liveroom.LiveRoom, error) {
	parse := gjson.Parse(content)
	streamInfo := parse.Get("roomInfo.tLiveInfo.tLiveStreamInfo.vStreamInfo.value")
	var urls []string
	streamInfo.ForEach(func(key, value gjson.Result) bool {
		urlPart := value.Get("sStreamName").String() + "." + value.Get("sFlvUrlSuffix").String() + "?" + value.Get("sFlvAntiCode").String()
		urls = append(urls, value.Get("sFlvUrl").String()+"/"+urlPart)
		return true
	})
	liveLineUrl := parse.Get("roomProfile.liveLineUrl").String()

	liveUrlByteData, err := base64.StdEncoding.DecodeString(liveLineUrl)
	if err != nil {
		return nil, errors.New("未开播或直播间不存在")
	}
	liveUrl, err := live(liveUrlByteData)
	if err != nil {
		return nil, errors.New("未开播或直播间不存在")
	}
	liveUrl = strings.ReplaceAll("https:"+liveUrl, "hls", "flv")
	liveUrl = strings.ReplaceAll(liveUrl, "m3u8", "flv")
	liveUrl = strings.ReplaceAll(liveUrl, "&ctype=tars_mobile", "")
	return &liveroom.LiveRoom{
		LiveUrl: liveUrl,
	}, nil
}

func live(byteData []byte) (string, error) {
	liveUrl := string(byteData)
	strs := strings.Split(liveUrl, "?")
	if len(strs) <= 1 {
		return "", errors.New("未开播或直播间不存在")
	}
	r := strings.Split(strs[0], "/")
	reg := regexp.MustCompile(`.(flv|m3u8)`)
	s := reg.ReplaceAllString(r[len(r)-1], "")
	c := strings.SplitN(strs[1], "&", 4)
	c1 := []string{}
	for _, str := range c {
		if str != "" {
			c1 = append(c1, str)
		}
	}
	n := make(map[string]string)
	for _, str := range c1 {
		cs := strings.Split(str, "=")
		n[cs[0]] = cs[1]
	}
	fm, err := url.QueryUnescape(n["fm"])
	if err != nil {
		return "", err
	}
	u := codec.Base64Decode(fm)
	p := strings.Split(u, "_")[0]
	f := strconv.Itoa(int(time.Now().UnixNano()))
	l := n["wsTime"]
	t := "0"
	hs := []string{p, t, s, f, l}
	h := strings.Join(hs, "_")
	m := codec.CalcMD5(h)
	y := c1[len(c1)-1]
	url := fmt.Sprintf("%s?wsSecret=%s&wsTime=%s&u=%s&seqid=%s&%s", strs[0], m, l, t, f, y)
	return url, nil
}

// 获取直播流详细信息
func (h *HuYa) GetStreamInfo(roomId string) ([]liveroom.StreamInfo, error) {
	streamInfo := make([]liveroom.StreamInfo, 0)
	url := "https://www.huya.com/" + roomId
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile(`"vMultiStreamInfo":(.*?),"iWebDefaultBitRate"`)
	submatch := reg.FindStringSubmatch(string(result))
	if submatch == nil {
		return nil, errors.Wrap(err, "查询失败！")
	}
	content := submatch[1]
	rates := gjson.Parse(content).Array()
	for _, rate := range rates {
		streamInfo = append(streamInfo, liveroom.StreamInfo{
			DisplayName: rate.Get("sDisplayName").String(),
			BitRate:     rate.Get("iBitRate").Int(),
			Url:         "",
		})
	}
	h.log.InfoFields("GetStreamInfo", logger.Fields{"streamInfo": streamInfo})
	return streamInfo, nil
}

func GetType(url string) string {
	if strings.Contains(url, "m3u8") {
		return "hls"
	}
	return "flv"
}

func (h *HuYa) GetRoomInfo(roomId string) (liveroom.LiveRoomInfo, error) {
	roomInfo := liveroom.LiveRoomInfo{}
	url := "https://www.huya.com/" + roomId
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return roomInfo, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	response, err := h.httpClient.Do(request)
	if err != nil {
		return roomInfo, err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return roomInfo, err
	}
	roomInfo.Anchor = matchRoomInfo(string(result), Nick)
	roomInfo.Avatar = matchRoomInfo(string(result), Avatar)
	roomInfo.GameFullName = matchRoomInfo(string(result), GameFullName)
	roomInfo.Screenshot = matchRoomInfo(string(result), Screenshot)
	roomInfo.RoomName = matchRoomInfo(string(result), RoomName)
	if roomInfo.RoomName == "" {
		roomInfo.RoomName = matchRoomInfo(string(result), Introduction)
	}
	roomInfo.OnLineCount, err = strconv.Atoi(matchRoomInfo(string(result), ActivityCount))
	if err != nil {
		roomInfo.OnLineCount = 0
	}
	liveStatus := matchRoomInfo(string(result), LiveStatus)
	h.log.InfoFields("直播状态", logger.Fields{"liveStatus": liveStatus})
	switch liveStatus {
	case "off":
		roomInfo.LiveStatus = 0
	case "on":
		roomInfo.LiveStatus = 2
	case "replay":
		roomInfo.LiveStatus = 1
	}
	if roomInfo.OnLineCount > 0 {
		roomInfo.LiveStatus = 2
	}
	roomInfo.Platform = "huya"
	roomInfo.PlatformName = liveroom.GetPlatform(roomInfo.Platform)
	roomInfo.RoomId = roomId
	return roomInfo, nil
}

func matchRoomInfo(res, regStr string) string {
	reg := regexp.MustCompile(regStr)
	match := reg.FindStringSubmatch(res)
	if match == nil {
		return ""
	}
	return match[1]
}
