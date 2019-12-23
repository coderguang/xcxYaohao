package define

import (
	"strconv"
	"sync"

	"github.com/coderguang/GameEngine_go/sgtime"
)

var staticStr []string

func init() {
	staticStr = []string{"打开次数", "新用户", "请求次数", "验证码请求次数", "短信成功", "短信失败", "绑定人数", "取消绑定", "分享", "unknow", "", "", "", "", ""}
}

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

type StatisticsData struct {
	Lock            sync.RWMutex
	ID              int `gorm:"primary_key":AUTO_INCREMENT`
	Time            string
	TimesData       map[int]int `gorm:"-"`
	StrData         string
	OpenTimes       int
	NewUser         int
	RequireTimes    int
	RandomCodeTimes int
	SmsSuccess      int
	SmsFailed       int
	BindTimes       int
	CancelTimes     int
	ShareTimes      int
	UserSize        int
}

func (data *StatisticsData) Reset() {

	data.Lock.Lock()
	defer data.Lock.Unlock()
	data.TimesData = make(map[int]int)
	now := sgtime.New()
	data.Time = sgtime.YearString(now) + sgtime.MonthString(now) + sgtime.DayString(now)

}

func (data *StatisticsData) String() string {
	data.Lock.Lock()
	defer data.Lock.Unlock()
	str := "\n\n========="
	for k, v := range data.TimesData {
		str += "\n" + staticStr[k] + ":" + strconv.Itoa(v)
	}
	str += "\n 总用户数:" + strconv.Itoa(data.UserSize)
	return str
}

func (data *StatisticsData) ParseToDb() {
	for k, v := range data.TimesData {
		switch k {
		case StatisticOpenTimes:
			data.OpenTimes = v
		case StatisticNewOpenTimes:
			data.NewUser = v
		case StatisticRequireTimes:
			data.RequireTimes = v
		case StatisticRandomCodeSend:
			data.RandomCodeTimes = v
		case StatisticSmsSuccess:
			data.SmsSuccess = v
		case StatisticSmsFail:
			data.SmsFailed = v
		case StatisticBindTimes:
			data.BindTimes = v
		case StatisticBindCancel:
			data.CancelTimes = v
		case StatisticShareTime:
			data.ShareTimes = v
		}
	}
}
