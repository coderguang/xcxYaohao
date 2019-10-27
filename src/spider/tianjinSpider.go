package spider

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
)

func TianjinOldDataSpider(cmd []string) {
	unitIndex := "http://apply.xkctk.jtys.tj.gov.cn/apply/norm/unitQuery.html"
	persionIndex := "http://apply.xkctk.jtys.tj.gov.cn/apply/norm/personQuery.html"

	startDt, err := sgtime.ParseInLocation(sgtime.FORMAT_TIME_NORMAL, "2014-02-01 00:00:00")
	if err != nil {
		sglog.Error("parse startDt error,", err)
	}

	go StartSpiderEx(define.CITY_TIANJIN, startDt, unitIndex, false)
	go StartSpiderEx(define.CITY_TIANJIN, startDt, persionIndex, true)
}
