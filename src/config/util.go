package config

import (
	"github.com/coderguang/GameEngine_go/sgcfg"
	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgthread"
)

var (
	globalUtilCfg *UtilCfg
)

type UtilCfg struct {
	BindID   int    `json:"bindId"`
	LuckID   int    `json:"luckId"`
	Port     string `json:"port"`
	Receiver string `json:"receiver"`
	Sign     string `json:"sign"`
	Tcid     string `json:"tcid"`
	Tckey    string `json:"tckey"`
	UnluckID int    `json:"unluckId"`
}

func initUtilCfg() {
	cfgFile := sgcfg.GetServerCfgDir() + "util.json"
	globalUtilCfg = new(UtilCfg)
	err := sgcfg.ReadCfg(cfgFile, globalUtilCfg)
	if err != nil {
		sglog.Error("init spider config error,", err)
		sgthread.DelayExit(2)
	}
}

func GetUtilCfg() *UtilCfg {
	return globalUtilCfg
}

func GetTxSmsCfg() (string, string) {
	return globalUtilCfg.Tcid, globalUtilCfg.Tckey
}

func GetSign()string{
	return globalUtilCfg.Sign
}

func GetBindId()int{
	return globalUtilCfg.BindID
}

func GetLuckId()int{
	return globalUtilCfg.LuckID   
}

func GetUnLuckId()int{
	return globalUtilCfg.UnluckID 
}