package liveroom

const (
	AreaInfosKey = "AreaInfos"
)

type LiveRoom struct {
	Platform     string    `json:"platform"`     //直播平台
	PlatformName string    `json:"platformName"` //直播平台名称
	RoomId       string    `json:"roomId"`       //房间ID
	RoomName     string    `json:"roomName"`     //房间名称
	Anchor       string    `json:"anchor"`       //主播
	Avatar       string    `json:"avatar"`       //头像
	OnLineCount  int       `json:"onLineCount"`  //在线人数
	LiveUrl      string    `json:"liveUrl"`      //真实直播流
	Quality      []Quality `json:"quality"`      //直播流详细信息
	IsLive       bool      `json:"isLive"`       //是否在直播
	Screenshot   string    `json:"screenshot"`   //房间封面图
	GameFullName string    `json:"gameFullName"` //类别
}

type Quality struct {
	Name string `json:"name"` //画质名称
	Url  string `json:"url"`  //地址
	Type string `json:"type"` //视频类型
}

type ArtQuality struct {
	Default bool   `json:"default"` //默认画质
	Html    string `json:"html"`    //画质名称
	Url     string `json:"url"`     //地址
}

type StreamInfo struct {
	DisplayName string `json:"displayName"` //画质名称
	BitRate     int64  `json:"bitRate"`     //画质码率
	Url         string `json:"url"`         //地址
}

type LiveRoomInfo struct {
	Platform     string `json:"platform"`     //直播平台
	PlatformName string `json:"platformName"` //直播平台名称
	RoomId       string `json:"roomId"`       //房间ID
	RoomName     string `json:"roomName"`     //房间名称
	Anchor       string `json:"anchor"`       //主播
	Avatar       string `json:"avatar"`       //头像
	OnLineCount  int    `json:"onLineCount"`  //在线人数
	Screenshot   string `json:"screenshot"`   //房间封面图
	GameFullName string `json:"gameFullName"` //类别
	LiveStatus   int    `json:"liveStatus"`   //直播状态,0-off,2-on,1-replay
}

// 平台代码转换
func GetPlatform(platform string) string {
	str := ""
	switch platform {
	case "huya":
		str = "虎牙"
	case "douyu":
		str = "斗鱼"
	case "bilibili":
		str = "哔哩哔哩"
	}
	return str
}

// 分区详情
type AreaInfo struct {
	Platform  string `json:"platform"`
	AreaId    string `json:"areaId"`
	AreaName  string `json:"areaName"`
	AreaPic   string `json:"areaPic"`
	ShortName string `json:"shortName"`
	TypeName  string `json:"typeName"`
	AreaType  string `json:"areaType"`
}

type LiveRoomInfoArray []LiveRoomInfo

func (array LiveRoomInfoArray) Len() int {
	return len(array)
}

func (array LiveRoomInfoArray) Less(i, j int) bool {
	return array[i].LiveStatus > array[j].LiveStatus //从小到大， 若为大于号，则从大到小
}

func (array LiveRoomInfoArray) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}
