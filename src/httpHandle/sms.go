package httpHandle

import (
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"
)

func sendConfirmMsg(phone string, randomCode string) YaoHaoNoticeError {
	return YAOHAO_OK
}

func sendCommonSms() {
	data.AddStatistic(define.StatisticSmsSuccess, 1)
	data.AddStatistic(define.StatisticSmsFail, 1)
}
