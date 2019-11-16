package spider

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sglog"
)

func HainanOldDataSpider(cmd []string) {
	unitIndex := "https://apply.hnjdctk.gov.cn/apply/app/status/norm/unit"
	persionIndex := "https://apply.hnjdctk.gov.cn/apply/app/status/norm/person"

	startDt, err := sgtime.ParseInLocation(sgtime.FORMAT_TIME_NORMAL, "2018-08-01 00:00:00")
	if err != nil {
		sglog.Error("parse startDt error,", err)
	}

	go StartSpiderEx(define.CITY_HAINAN, startDt, unitIndex, false)
	go StartSpiderEx(define.CITY_HAINAN, startDt, persionIndex, true)

}
