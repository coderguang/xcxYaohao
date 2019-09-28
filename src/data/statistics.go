package data

import (
	"strconv"
	"sync"

	"github.com/coderguang/GameEngine_go/sgtime"
)

const (
	StatisticOpenTimes      int = 0 //require op=time
	StatisticNewOpenTimes   int = 1
	StatisticRequireTimes   int = 2 // op = search
	StatisticRandomCodeSend int = 3
	StatisticSmsSuccess     int = 4
	StatisticSmsFail        int = 5
	StatisticBindTimes      int = 6
	StatisticBindCancel     int = 7
	StatisticShareTime      int = 8
	StatisticMax            int = 9
)

var (
	globalStatics *StatisticsData
	staticStr     []string
)

func init() {
	globalStatics = new(StatisticsData)
	globalStatics.Reset()
	staticStr = []string{"打开次数", "新用户", "请求次数", "验证码请求次数", "短信成功", "短信失败", "绑定人数", "取消绑定", "分享", "unknow", "", "", "", "", ""}
}

type StatisticsData struct {
	Lock      sync.RWMutex
	Time      string `gorm:"primary_key"`
	TimesData map[int]int
}

func (data *StatisticsData) Reset() {

	data.Lock.Lock()
	defer data.Lock.Unlock()
	data.TimesData = make(map[int]int)
	now := sgtime.New()
	data.Time = sgtime.YearString(now) + sgtime.MonthString(now)

}

func (data *StatisticsData) String() string {
	data.Lock.Lock()
	defer data.Lock.Unlock()
	str := "\n\n========="
	for k, v := range data.TimesData {
		str += "\n" + staticStr[k] + ":" + strconv.Itoa(v)
	}
	return str
}

func AddStatistic(logType int, addTimes int) {
	globalStatics.Lock.Lock()
	defer globalStatics.Lock.Unlock()
	if v, ok := globalStatics.TimesData[logType]; ok {
		globalStatics.TimesData[logType] = v + addTimes
	} else {
		globalStatics.TimesData[logType] = addTimes
	}
}
