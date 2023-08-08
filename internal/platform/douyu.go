package platform

import (
	"changeme/internal/global"
	"changeme/internal/liveroom"
	"changeme/pkg/request"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/patrickmn/go-cache"
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
	cache      cache.Cache
}

type RealUrlResp struct {
	Code int64 `json:"code"`
	Data Data  `json:"data"`
}

type Data struct {
	Settings        []Setting `json:"settings"`
	URL             string    `json:"url"`
	Rate            int64     `json:"rate"`
	Pass            int64     `json:"pass"`
	ShareOffsetTime int64     `json:"share_offset_time"`
}

type Setting struct {
	Name    string `json:"name"`
	Rate    int64  `json:"rate"`
	HighBit int64  `json:"high_bit"`
}

func NewDoYu() DouYu {
	return DouYu{
		httpClient: http.Client{},
		log:        logger.NewCustomLogger(Douyu),
		cache:      *global.Cache,
	}
}

func md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func (d *DouYu) getDid() (string, error) {
	timeStamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	url := "https://passport.douyu.com/lapi/did/api/get?client_id=25&_=" + timeStamp + "&callback=axiosJsonpCallback1"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	req.Header.Set("referer", "https://m.douyu.com/")
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	re := regexp.MustCompile(`axiosJsonpCallback1\((.*)\)`)
	match := re.FindStringSubmatch(string(body))
	var result map[string]map[string]string
	json.Unmarshal([]byte(match[1]), &result)
	return result["data"]["did"], nil
}

func (d *DouYu) GetRealUrl(roomId, streamType string) (*liveroom.LiveRoom, error) {
	did, err := d.getDid()
	if err != nil {
		return nil, err
	}
	var timestamp = time.Now().Unix()
	liveurl := "https://m.douyu.com/" + roomId
	client := &http.Client{}
	r, _ := http.NewRequest("GET", liveurl, nil)
	r.Header.Add("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1")
	r.Header.Add("upgrade-insecure-requests", "1")
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	roomidreg := regexp.MustCompile(`(?i)rid":(\d{1,8}),"vipId`)
	roomidres := roomidreg.FindStringSubmatch(string(body))
	if roomidres == nil {
		return nil, errors.New("roomid not found")
	}
	realroomid := roomidres[1]
	reg := regexp.MustCompile(`(?i)(function ub98484234.*)\s(var.*)`)
	res := reg.FindStringSubmatch(string(body))
	nreg := regexp.MustCompile(`(?i)eval.*;}`)
	strfn := nreg.ReplaceAllString(res[0], "strc;}")
	vm := goja.New()
	_, err = vm.RunString(strfn)
	if err != nil {
		return nil, err
	}
	jsfn, ok := goja.AssertFunction(vm.Get("ub98484234"))
	if !ok {
		return nil, errors.New("ub98484234 not found")
	}
	result, err := jsfn(
		goja.Undefined(),
		vm.ToValue("ub98484234"),
	)
	if err != nil {
		return nil, err
	}
	nres := fmt.Sprintf("%s", result)
	nnreg := regexp.MustCompile(`(?i)v=(\d+)`)
	nnres := nnreg.FindStringSubmatch(nres)
	unrb := fmt.Sprintf("%v%v%v%v", realroomid, did, timestamp, nnres[1])
	rb := md5V3(unrb)
	nnnreg := regexp.MustCompile(`(?i)return rt;}\);?`)
	strfn2 := nnnreg.ReplaceAllString(nres, "return rt;}")
	strfn3 := strings.Replace(strfn2, `(function (`, `function sign(`, -1)
	strfn4 := strings.Replace(strfn3, `CryptoJS.MD5(cb).toString()`, `"`+rb+`"`, -1)

	_, err = vm.RunString(strfn4)
	if err != nil {
		return nil, err
	}
	sign, ok := goja.AssertFunction(vm.Get("sign"))
	if !ok {
		return nil, errors.New("sign not found")
	}
	param, err := sign(
		goja.Undefined(),
		vm.ToValue(realroomid),
		vm.ToValue(did),
		vm.ToValue(timestamp),
	)
	if err != nil {
		return nil, err
	}
	params := fmt.Sprintf("%s&ver=22107261&rid=%s&rate=-1", param, realroomid)
	r1, n4err := http.Post("https://m.douyu.com/api/room/ratestream", "application/x-www-form-urlencoded", strings.NewReader(params))
	if n4err != nil {
		panic(n4err)
	}
	defer r1.Body.Close()
	body1, _ := io.ReadAll(r1.Body)
	var realUrlResp RealUrlResp
	json.Unmarshal(body1, &realUrlResp)
	var hlsUrl string
	if realUrlResp.Code != 0 {
		return nil, errors.New("roomid not found")
	}
	hlsUrl = strings.Replace(realUrlResp.Data.URL, "http://", "https://", 1)
	var realUrl string
	switch streamType {
	case "hls":
		realUrl = hlsUrl
	case "flv":
		realUrl = strings.Replace(hlsUrl, "m3u8", "flv", 1)
	case "xs":
		realUrl = strings.Replace(hlsUrl, "m3u8", "xs", 1)
	}
	room := new(liveroom.LiveRoom)
	room.LiveUrl = realUrl
	room.Platform = Platform
	room.PlatformName = liveroom.GetPlatform(room.Platform)
	room.RoomId = roomId
	return room, nil
}

func (d *DouYu) GetLiveUrl(roomId string) (*liveroom.LiveRoom, error) {
	return d.GetRealUrl(roomId, "hls")
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

	roomInfo = liveroom.LiveRoomInfo{
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
	}
	if _, ok := global.FocusMap.Get(global.FormatKey(liveroom.FocusKey, Douyu, roomId)); ok {
		roomInfo.Favorite = true
	}
	return roomInfo, nil
}

// 斗鱼获取推荐信息
func (d *DouYu) GetRecommend(page, pageSize int) ([]liveroom.LiveRoomInfo, error) {
	if list, ok := global.Cache.Get(global.FormatKey(liveroom.RecommendKey, Douyu, strconv.Itoa(page), strconv.Itoa(pageSize))); ok {
		return list.([]liveroom.LiveRoomInfo), nil
	}
	start := pageSize*(page-1)/8 + lo.If(pageSize*(page-1)/8 == 0, 0).Else(1)
	start = lo.If(start == 0, 1).Else(start)
	startIndex := pageSize * (page - 1) % 8
	end := pageSize*page/8 + lo.If(pageSize*page%8 == 0, 0).Else(1)
	endIndex := pageSize * page % 8

	roomInfos := []liveroom.LiveRoomInfo{}

	d.log.InfoFields("End", logger.Fields{"End": end})

	for i := start; i <= end; i++ {

		var resp struct {
			Code int `json:"code"`
			Data struct {
				List []struct {
					RoomId     int    `json:"rid"`      //房间ID
					RoomName   string `json:"roomName"` //房间名称
					Anchor     string `json:"nickname"` //主播
					Avatar     string `json:"avatar"`   //头像
					Screenshot string `json:"roomSrc"`  //房间封面图
				} `json:"list"`
			} `json:"data"`
		}
		err := request.HTTP().GET(fmt.Sprintf("https://m.douyu.com/api/room/list?page=%d&type=", i)).
			BindJSON(&resp).Do()
		if err != nil {
			d.log.Error(err.Error())
			return nil, err
		}
		list := []liveroom.LiveRoomInfo{}
		if resp.Code == 0 {
			for _, res := range resp.Data.List {
				room, err := d.GetRoomInfo(strconv.Itoa(res.RoomId))
				if err != nil {
					d.log.Error(err.Error())
					return nil, err
				}
				list = append(list, room)
			}
			roomInfos = append(roomInfos, list...)
		}
	}
	d.log.InfoFields("roomInfos", logger.Fields{"roomInfos length": len(roomInfos)})
	roomInfos = roomInfos[startIndex : len(roomInfos)-endIndex]
	global.Cache.Set(global.FormatKey(liveroom.RecommendKey, Douyu, strconv.Itoa(page), strconv.Itoa(pageSize)), roomInfos, 10*time.Minute)
	return roomInfos, nil
}
