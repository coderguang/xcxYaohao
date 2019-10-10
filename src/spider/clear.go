package spider

import (
	"time"
	"xcxYaohao/src/data"
	"xcxYaohao/src/db"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgthread"
	"github.com/coderguang/GameEngine_go/sgtime"
)

func InitClear() {

	for {

		nowTime := time.Now()
		normalTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 0, nowTime.Location())
		timeInt := normalTime.Sub(nowTime)
		sleepTime := int(timeInt/time.Second) + 10 // +1 for avoid loop run in that second time
		sglog.Info("next clear timer will run after ", sleepTime, " seconds in ", sleepTime, sgtime.NormalString(sgtime.TransfromTimeToDateTime(normalTime)))
		//sleepTime = 60
		sgthread.SleepBySecond(sleepTime)

		sglog.Info("start to run clear openid data")

		dbData := data.ClearData()
		if dbData != nil {
			dbData.ParseToDb()
			db.UpdateStatisData(dbData)
		} else {
			sglog.Error("can't get statistic data")
		}
	}
}
