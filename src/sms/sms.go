package sms

import (
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtc/tcsms"

	"github.com/coderguang/GameEngine_go/sglog"
)

func SendRandomCode(phone string, title string, code string, randomCode string) error {

	cityName := config.GetCityName(title)

	code = "***" + code[len(code)-4:len(code)]

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetBindId(), cityName, code, randomCode)

	if err != nil {
		sglog.Error("send sms randomCode error", err)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return err
	} else {
		sglog.Info("recv sms randmcode,", res)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return nil
	}
}

func SendLuck(phone string, title string, time string) error {

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetLuckId(), cityName, time)

	if err != nil {
		sglog.Error("send sms luck error", err)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return err
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return err
	}
}

func SendUnLuck(phone string, title string, time string) error {

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetUnLuckId(), cityName, time)

	if err != nil {
		sglog.Error("send sms luck error", err)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return err
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return err
	}
}
