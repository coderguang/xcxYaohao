package sms

import (
	"errors"
	"xcxYaohao/src/config"
	"xcxYaohao/src/data"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtc/tcsms"

	"github.com/coderguang/GameEngine_go/sglog"
)

var (
	globalSmsFlag bool
)

func init() {
	globalSmsFlag = true
}

func SendRandomCode(phone string, title string, code string, randomCode string) error {

	if !globalSmsFlag {
		return errors.New("sms flag false,would not send sms")
	}

	ignorePhoneList := config.GetUtilCfg().IgnorePhone
	for _, v := range ignorePhoneList {
		if phone == v {
			sglog.Debug("ignore phone num:", v)
			return nil
		}
	}

	cityName := config.GetCityName(title)

	code = "***" + code[len(code)-4:len(code)]

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetBindId(), cityName, code, randomCode)

	if err != nil {
		sglog.Error("send sms randomCode error", err)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return err
	} else {
		sglog.Info("recv sms randmcode,", res)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return nil
	}
}

func SendLuck(phone string, title string, time string) error {

	if !globalSmsFlag {
		return errors.New("sms flag false,would not send sms")
	}

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetLuckId(), cityName, time)

	if err != nil {
		sglog.Error("send sms luck error", err)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return err
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return err
	}
}

func SendUnLuck(phone string, title string, time string) error {

	if !globalSmsFlag {
		return errors.New("sms flag false,would not send sms")
	}

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetUnLuckId(), cityName, time)

	if err != nil {
		sglog.Error("send sms unluck error", err)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return err
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return err
	}
}

func SendUnLuckAndTimeout(phone string, title string, time string) error {

	if !globalSmsFlag {
		return errors.New("sms flag false,would not send sms")
	}

	cityName := config.GetCityName(title)

	// 发送短信
	res, err := tcsms.SingleSend(config.GetSign(), 86, phone, config.GetTimeOutId(), cityName, time)

	if err != nil {
		sglog.Error("send sms timeout unluck error", err)
		data.AddStatistic(define.StatisticSmsFail, 1)
		return err
	} else {
		sglog.Info("recv sms luck,", res)
		data.AddStatistic(define.StatisticSmsSuccess, 1)
		return err
	}
}
