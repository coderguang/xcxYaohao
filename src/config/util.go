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
	BindID      int      `json:"bindId"`
	LuckID      int      `json:"luckId"`
	Port        string   `json:"port"`
	Receiver    string   `json:"receiver"`
	Sign        string   `json:"sign"`
	UnluckID    int      `json:"unluckId"`
	NoCache     bool     `json:"noCache"`
	TimeOutId   int      `json:"timeoutId"`
	IgnorePhone []string `json:"ignorePhone`
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

func GetSign() string {
	return globalUtilCfg.Sign
}

func GetBindId() int {
	return globalUtilCfg.BindID
}

func GetLuckId() int {
	return globalUtilCfg.LuckID
}

func GetUnLuckId() int {
	return globalUtilCfg.UnluckID
}

func GetTimeOutId() int {
	return globalUtilCfg.TimeOutId
}

func IsNoCache() bool {
	return globalUtilCfg.NoCache
}
