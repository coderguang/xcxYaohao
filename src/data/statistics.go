package data

import "xcxYaohao/src/define"

var (
	globalStatics *define.StatisticsData
)

func init() {
	globalStatics = new(define.StatisticsData)
	globalStatics.Reset()
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
