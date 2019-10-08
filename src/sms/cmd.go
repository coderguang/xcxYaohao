package sms

import (
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
)

var testPhone string = "15711839048"

func SendTestRandom(cmd []string) {
	//SendRandomCode(testPhone, "guanzhou", "1234567890123", "1311")
	SendRandomCode(testPhone, define.CITY_GUANGZHOU, "1234", "1311")
}

func SendTestLuck(cmd []string) {
	SendLuck(testPhone, define.CITY_SHENZHEN, "201909")
}

func SendTestUnLuck(cmd []string) {
	SendUnLuck(testPhone, define.CITY_HANGZHOU, "201912")
}

func ShowCurrentSmsFlang(cmd []string) {
	sglog.Info("current sms flag:", globalSmsFlag)
}
func ChangeCurrentSmsFlang(cmd []string) {
	globalSmsFlag = !globalSmsFlag
	sglog.Info("after change current sms flag:", globalSmsFlag)
}
