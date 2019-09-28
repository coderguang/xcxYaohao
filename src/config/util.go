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
	Port     string `json:"port"`
	Receiver string `json:"receiver"`
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
