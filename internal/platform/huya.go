package platform

import (
	"changeme/internal/global"
	"changeme/internal/liveroom"
	"changeme/pkg/request"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/guonaihong/gout"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

type HuYa struct {
	httpClient http.Client
	log        *wails.CustomLogger
	cache      cache.Cache
}

const (
	Nick          = `\"nick\":\"(.*?)\",`
	RoomName      = `\"roomName\":\"(.*?)\",`
	Screenshot    = `\"screenshot\":\"(.*?)\",`
	Avatar        = `\"avatar180\":\"(.*?)\",`
	GameFullName  = `\"gameFullName\":\"(.*?)\",`
	ActivityCount = `\"activityCount\":(.*?),`
	Introduction  = `\"introduction\":\"(.*?)\",`
	LiveStatus    = `\"liveStatus-(.*?)\"`
)

var (
	HuyaTopAreaMap = map[string]string{
		"1": "网游竞技",
		"2": "单机热游",
		"3": "手游休闲",
		"8": "娱乐天地",
	}
)

func NewHuYa() HuYa {
	huya := HuYa{
		httpClient: http.Client{Timeout: time.Second * 5},
		log:        logger.NewCustomLogger("huya"),
		cache:      *global.Cache,
	}
	huya.InitAreaCache()
	return huya
}

// 获取真实直播流
func (h *HuYa) GetLiveUrl(roomId string) (*liveroom.LiveRoom, error) {
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
	result, err := io.ReadAll(response.Body)
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
	streamInfo.ForEach(func(_, value gjson.Result) bool {
		urlStr := fmt.Sprintf("%s/%s.%s?%s",
			value.Get("sFlvUrl").String(),
			value.Get("sStreamName").String(),
			value.Get("sFlvUrlSuffix").String(),
			parseAntiCode(value.Get("sHlsAntiCode").String(), getAnonymousUid(), value.Get("sStreamName").String()))
		urls = append(urls, urlStr)
		return true
	})
	liveUrl := urls[rand.Intn(len(urls)-1)]
	return &liveroom.LiveRoom{
		LiveUrl: strings.Replace(liveUrl, "http://", "https://", 1),
	}, nil
}

func parseAntiCode(anticode, uid, streamName string) string {
	qr, err := url.ParseQuery(anticode)
	if err != nil {
		return ""
	}
	uidInt, _ := strconv.Atoi(uid)
	qr.Set("ver", "1")
	qr.Set("sv", "2110211124")
	qr.Set("seqid", strconv.FormatInt(time.Now().Unix()*1000+int64(uidInt), 10))
	qr.Set("uid", uid)
	qr.Set("uuid", strconv.Itoa(getUuid()))
	ss := MD5([]byte(fmt.Sprintf("%s|%s|%s", qr.Get("seqid"), qr.Get("ctype"), qr.Get("t"))))

	decodeString, _ := base64.StdEncoding.DecodeString(qr.Get("fm"))
	fm := string(decodeString)
	fm = strings.ReplaceAll(fm, "$0", qr.Get("uid"))
	fm = strings.ReplaceAll(fm, "$1", streamName)
	fm = strings.ReplaceAll(fm, "$2", ss)
	fm = strings.ReplaceAll(fm, "$3", qr.Get("wsTime"))

	qr.Del("fm")
	qr.Set("wsSecret", MD5([]byte(fm)))
	if qr.Has("txyp") {
		qr.Del("txyp")
	}
	return qr.Encode()
}

func getAnonymousUid() string {
	urlStr := "https://udblgn.huya.com/web/anonymousLogin"
	body := "{\n        \"appId\": 5002,\n        \"byPass\": 3,\n        \"context\": \"\",\n        \"version\": \"2.4\",\n        \"data\": {}\n    }"
	resp, err := http.Post(urlStr, "application/json", strings.NewReader(body))
	if err != nil {
		return ""
	}
	result, _ := io.ReadAll(resp.Body)
	return gjson.Parse(string(result)).Get("data.uid").String()
}

func getUuid() int {
	now := time.Now().Unix()
	random := int64(rand.Intn(1000))
	return int((now%10000000000*1000 + random) % 4294967295)
}

func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
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
	result, err := io.ReadAll(response.Body)
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
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return roomInfo, err
	}
	//h.log.InfoFields("房间详情元数据", logger.Fields{"result": string(result)})
	roomInfo.Anchor = matchString(string(result), Nick)
	roomInfo.Avatar = matchString(string(result), Avatar)
	if roomInfo.Avatar != "" && !strings.Contains(roomInfo.Avatar, "https") {
		roomInfo.Avatar = strings.ReplaceAll(roomInfo.Avatar, "http", "https")
	}
	roomInfo.GameFullName = matchString(string(result), GameFullName)
	roomInfo.Screenshot = matchString(string(result), Screenshot)
	if roomInfo.Screenshot != "" && !strings.Contains(roomInfo.Screenshot, "https") {
		roomInfo.Screenshot = strings.ReplaceAll(roomInfo.Screenshot, "http", "https")
	}
	roomInfo.RoomName = matchString(string(result), RoomName)
	if roomInfo.RoomName == "" {
		roomInfo.RoomName = matchString(string(result), Introduction)
	}
	roomInfo.OnLineCount, err = strconv.Atoi(matchString(string(result), ActivityCount))
	if err != nil {
		roomInfo.OnLineCount = 0
	}
	liveStatus := matchString(string(result), LiveStatus)
	h.log.InfoFields("直播状态", logger.Fields{"liveStatus": liveStatus})
	switch liveStatus {
	case "off":
		roomInfo.LiveStatus = 0
	case "on":
		roomInfo.LiveStatus = 2
	case "on on-match":
		roomInfo.LiveStatus = 2
	case "replay":
		roomInfo.LiveStatus = 1
	}
	roomInfo.Platform = "huya"
	roomInfo.PlatformName = liveroom.GetPlatform(roomInfo.Platform)
	roomInfo.RoomId = roomId
	if _, ok := global.FocusMap.Get(global.FormatKey(liveroom.FocusKey, roomInfo.Platform, roomInfo.RoomId)); ok {
		roomInfo.Favorite = true
	}
	return roomInfo, nil
}

func matchString(res, regStr string) string {
	reg := regexp.MustCompile(regStr)
	match := reg.FindStringSubmatch(res)
	if match == nil {
		return ""
	}
	return match[1]
}

// 初始化分区
func (h *HuYa) InitAreaCache() []liveroom.AreaInfo {
	areaInfos := []liveroom.AreaInfo{}
	for key, val := range HuyaTopAreaMap {
		infos, err := h.GetSimgleArea(key, val)
		if err != nil {
			break
		}
		areaInfos = append(areaInfos, infos...)
	}
	h.cache.Set(global.FormatKey(liveroom.AreaInfosKey, Huya), areaInfos, 30*time.Minute)
	return areaInfos
}

// 获取单个分区
func (h *HuYa) GetSimgleArea(areaCode, typeName string) ([]liveroom.AreaInfo, error) {
	areaInfos := []liveroom.AreaInfo{}
	var resp struct {
		GameList []struct {
			Gid          int    `json:"gid"`
			GameFullName string `json:"gameFullName"`
		} `json:"gameList"`
	}
	//var resp string
	err := request.HTTP().GET(fmt.Sprintf("https://m.huya.com/cache.php?m=Game&do=ajaxGameList&bussType=%s", areaCode)).
		SetHeader(gout.H{
			"User-Agent":   "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36",
			"Content-Type": "application/x-www-form-urlencoded",
		}).BindJSON(&resp).Do()
	if err != nil {
		return nil, err
	}
	h.log.InfoFields("areaInfos", logger.Fields{"areaInfos": resp})
	for _, game := range resp.GameList {
		areaInfos = append(areaInfos, liveroom.AreaInfo{
			Platform:  Huya,
			AreaId:    strconv.Itoa(game.Gid),
			AreaName:  game.GameFullName,
			AreaPic:   "https://huyaimg.msstatic.com/cdnimage/game/" + strconv.Itoa(game.Gid) + "-MS.jpg",
			ShortName: "",
			TypeName:  typeName,
			AreaType:  areaCode,
		})
	}
	return areaInfos, nil
}

// 获取所有分区信息
func (h *HuYa) GetAllAreaInfo() ([]liveroom.AreaInfo, error) {
	if infos, ok := h.cache.Get(global.FormatKey(liveroom.AreaInfosKey, Huya)); ok {
		return infos.([]liveroom.AreaInfo), nil
	}
	return h.InitAreaCache(), nil
}

// 获取虎牙推荐
func (h *HuYa) GetRecommend(page, pageSize int) ([]liveroom.LiveRoomInfo, error) {
	if list, ok := global.Cache.Get(global.FormatKey(liveroom.RecommendKey, Huya, strconv.Itoa(page), strconv.Itoa(pageSize))); ok {
		return list.([]liveroom.LiveRoomInfo), nil
	}
	realPage := page/6 + 1
	start := (page - 1) * pageSize % 120
	if pageSize == 10 {
		realPage = page/12 + 1
		start = (page - 1) * pageSize % 120
	}
	var resp struct {
		Status int `json:"status"`
		Data   struct {
			Datas []struct {
				ProfileRoom  string `json:"profileRoom"`
				Gid          string `json:"gid"`
				GameFullName string `json:"gameFullName"`
				RoomName     string `json:"roomName"`
				Nick         string `json:"nick"`
				Screenshot   string `json:"screenshot"`
				Avatar       string `json:"avatar180"`
				TatalCount   string `json:"totalCount"`
			} `json:"datas"`
		} `json:"data"`
	}
	//var resp string
	err := request.HTTP().GET(fmt.Sprintf("https://www.huya.com/cache.php?m=LiveList&do=getLiveListByPage&tagAll=0&page=%d", realPage)).
		BindJSON(&resp).Do()
	if err != nil {
		return nil, err
	}
	h.log.InfoFields("roomInfos", logger.Fields{"roomInfos": resp, "start": start, "size": start + pageSize})
	roomInfos := make([]liveroom.LiveRoomInfo, 0, start+pageSize)
	if resp.Status == 200 {
		for i := start; i < start+pageSize; i++ {
			if i >= len(resp.Data.Datas) {
				break
			}
			count, err := strconv.Atoi(resp.Data.Datas[i].TatalCount)
			if err != nil {
				h.log.InfoFields("GetRoomInfo Count Err", logger.Fields{"count": count})
				break
			}
			roomInfo := liveroom.LiveRoomInfo{
				Platform:     Huya,
				PlatformName: liveroom.GetPlatform(Huya),
				RoomId:       resp.Data.Datas[i].ProfileRoom,
				RoomName:     resp.Data.Datas[i].RoomName,
				Anchor:       resp.Data.Datas[i].Nick,
				Avatar:       lo.If(resp.Data.Datas[i].Avatar != "" && !strings.Contains(resp.Data.Datas[i].Avatar, "https"), strings.ReplaceAll(resp.Data.Datas[i].Avatar, "http", "https")).Else(resp.Data.Datas[i].Avatar),
				OnLineCount:  count,
				Screenshot:   lo.If(resp.Data.Datas[i].Screenshot != "" && !strings.Contains(resp.Data.Datas[i].Screenshot, "https"), strings.ReplaceAll(resp.Data.Datas[i].Screenshot, "http", "https")).Else(resp.Data.Datas[i].Screenshot),
				GameFullName: resp.Data.Datas[i].GameFullName,
				LiveStatus:   2,
			}
			if _, ok := global.FocusMap.Get(global.FormatKey(liveroom.FocusKey, roomInfo.Platform, roomInfo.RoomId)); ok {
				roomInfo.Favorite = true
			}
			roomInfos = append(roomInfos, roomInfo)
		}
	}
	global.Cache.Set(global.FormatKey(liveroom.RecommendKey, Huya, strconv.Itoa(page), strconv.Itoa(pageSize)), roomInfos, 10*time.Minute)
	return roomInfos, nil
}
