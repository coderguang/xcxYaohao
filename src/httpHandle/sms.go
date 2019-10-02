package httpHandle

import (
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/zboyco/gosms"
)

func SendRandomCode(phone string, title string, code string, randomCode string) YaoHaoNoticeError {
	smsId, smsKey := config.GetTxSmsCfg()
	// 创建Sender
	sender := &gosms.QSender{
		AppID:  smsId,  // appid
		AppKey: smsKey, // appkey
	}

	cityName := config.GetCityName(title)

	code = "***" + code[len(code)-4:len(code)]

	// 发送短信
	res, err := sender.SingleSend(
		config.GetSign(),   // 短信签名，此处应填写审核通过的签名内容，非签名 ID，如果使用默认签名，该字段填 ""
		86,                 // 国家号
		phone,              // 手机号
		config.GetBindId(), // 短信正文ID
		cityName,
		code,
		randomCode, // 参数1
	)
	if err != nil {
		sglog.Error("send sms randomCode error", err)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return YAOHAO_ERR_SMS_RESULT_PARSE_ERROR
	} else {
		sglog.Info("recv sms randmcode,", res)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return YAOHAO_OK
	}
}

func SendLuck(phone string, title string, time string) YaoHaoNoticeError {
	smsId, smsKey := config.GetTxSmsCfg()
	// 创建Sender
	sender := &gosms.QSender{
		AppID:  smsId,  // appid
		AppKey: smsKey, // appkey
	}

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := sender.SingleSend(
		config.GetSign(),   // 短信签名，此处应填写审核通过的签名内容，非签名 ID，如果使用默认签名，该字段填 ""
		86,                 // 国家号
		phone,              // 手机号
		config.GetLuckId(), // 短信正文ID
		cityName,           // 参数1
		time,
	)
	if err != nil {
		sglog.Error("send sms luck error", err)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return YAOHAO_ERR_SMS_RESULT_PARSE_ERROR
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return YAOHAO_OK
	}
}

func SendUnLuck(phone string, title string, time string) YaoHaoNoticeError {
	smsId, smsKey := config.GetTxSmsCfg()
	// 创建Sender
	sender := &gosms.QSender{
		AppID:  smsId,  // appid
		AppKey: smsKey, // appkey
	}

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := sender.SingleSend(
		config.GetSign(),     // 短信签名，此处应填写审核通过的签名内容，非签名 ID，如果使用默认签名，该字段填 ""
		86,                   // 国家号
		phone,                // 手机号
		config.GetUnLuckId(), // 短信正文ID
		cityName,             // 参数1
		time,
	)
	if err != nil {
		sglog.Error("send sms luck error", err)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return YAOHAO_ERR_SMS_RESULT_PARSE_ERROR
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return YAOHAO_OK
	}
}
