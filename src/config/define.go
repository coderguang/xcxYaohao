package config

import (
	"sync"

	"github.com/coderguang/GameEngine_go/sgthread"

	"github.com/coderguang/GameEngine_go/sglog"

	"github.com/coderguang/GameEngine_go/sgcfg"
)

var (
	globalCfgs *SecureSpiderCfg
)

func init() {
	globalCfgs = new(SecureSpiderCfg)
	globalCfgs.Data = make(map[string]SpiderCfg)
}

type SpiderCfg struct {
	Title      string   `json:"title"`
	IndexURL   string   `json:"indexUrl"`
	AllowUrls  []string `json:"allowUrls"`
	IgnoreUrls []string `json:"ignoreUrls"`
	ResultDate int      `json:"resultDate"`
}

type SecureSpiderCfg struct {
	Data map[string]SpiderCfg
	Lock sync.RWMutex
}

func InitCfg() {
	cfgFile := sgcfg.GetServerCfgDir() + "spider.json"
	spiderCfgs := []SpiderCfg{}
	err := sgcfg.ReadCfg(cfgFile, &spiderCfgs)
	if err != nil {
		sglog.Error("init spider config error,", err)
		sgthread.DelayExit(2)
	}

	globalCfgs.Lock.Lock()
	defer globalCfgs.Lock.Unlock()

	for _, v := range spiderCfgs {
		globalCfgs.Data[v.Title] = v
	}

}
