package httpHandle

import (
	"xcxYaohao/src/data"

	"github.com/coderguang/GameEngine_go/sglog"
)

func NoticeCurrentMonthDataUpdate(title string, curTime string) {
	finalNotice := data.GetNoticeFinalTime(title)
	if finalNotice >= curTime {
		sglog.Info("notice already ", title, curTime, ",last notice is ", finalNotice)
		return
	}

}
