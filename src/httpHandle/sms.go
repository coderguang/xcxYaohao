package httpHandle

import "xcxYaohao/src/data"

func sendConfirmMsg(phone string, randomCode string) YaoHaoNoticeError {
	return YAOHAO_OK
}

func sendCommonSms() {
	data.AddStatistic(data.StatisticSmsSuccess, 1)
	data.AddStatistic(data.StatisticSmsFail, 1)
}
