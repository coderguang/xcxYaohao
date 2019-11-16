package config

import (
	"errors"
	"sync"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sgtc/tcsms"

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
	HTTP       string   `json:"http"`
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

	sglog.Info("load spider config ok,size=", len(globalCfgs.Data))

	initUtilCfg()

	tcCfg := sgcfg.GetServerCfgDir() + "tc_sms.json"
	err = tcsms.InitTcSms(tcCfg)
	if err != nil {
		sglog.Error("tc sms config error,err:", err)
		sgthread.DelayExit(2)
	}
	sglog.Info("tc sms init ok")
}

func GetSpiderCfg(title string) (SpiderCfg, error) {
	globalCfgs.Lock.Lock()
	defer globalCfgs.Lock.Unlock()

	if v, ok := globalCfgs.Data[title]; ok {
		return v, nil
	}
	return SpiderCfg{}, errors.New("no this title spider config,title:" + title)
}

func GetTitleList() []string {
	titlelist := []string{}
	for k := range globalCfgs.Data {
		titlelist = append(titlelist, k)
	}
	return titlelist
}

func IsSupportCity(title string) bool {
	if _, err := GetSpiderCfg(title); err != nil {
		if title == define.CITY_HAINAN {
			return true
		}
		return false
	}
	return true
}
