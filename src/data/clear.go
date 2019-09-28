package data

import (
	"time"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/coderguang/GameEngine_go/sgtime"
)

func InitClear() {

	for {

		nowTime := time.Now()
		normalTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 0, nowTime.Location())
		timeInt := normalTime.Sub(nowTime)
		sleepTime := int(timeInt/time.Second) + 1 // +1 for avoid loop run in that second time
		sglog.Info("next clear timer will run after ", sleepTime, " seconds in ", sleepTime, sgtime.NormalString(sgtime.TransfromTimeToDateTime(normalTime)))
		sleepTime = 60
		sgthread.SleepBySecond(sleepTime)

		sglog.Info("start to run clear openid data")
		globalOpenIds.Lock.Lock()
		now := sgtime.New()
		idset := make(map[string]bool)
		for k, v := range globalOpenIds.Data {
			if _, ok := idset[v.Openid]; !ok {
				idset[v.Openid] = true
			}
			if sgtime.GetTotalSecond(now)-sgtime.GetTotalSecond(sgtime.TransfromTimeToDateTime(*v.Time)) > 3600 {
				sglog.Debug("delete openid:", v.Openid, ",code:", v.Code)
				delete(globalOpenIds.Data, k)
			}
		}
		globalOpenIds.Lock.Unlock()
		sglog.Info("clear openid data complete,user size:", len(idset), globalStatics)

		//sgmail.SendMail("xcxYaohao statistic", []string{config.GetUtilCfg().Receiver}, globalStatics.String())

		globalStatics.Reset()

	}
}
