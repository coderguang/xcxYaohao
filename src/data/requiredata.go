package data

import (
	"errors"
	"xcxYaohao/src/define"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgthread"

	"github.com/coderguang/GameEngine_go/sgcfg"

	"github.com/coderguang/GameEngine_go/sgwx/sgwxdef"

	"github.com/coderguang/GameEngine_go/sgwx/sgwxopenid"

	"github.com/coderguang/GameEngine_go/sgali/alisms"
)

var (
	globalOpenIds     *define.SecureWxOpenid
	globalWxOpenIdCfg *sgwxdef.WxAppidCfg
	globalAliAppidCfg *alisms.AliAppidCfg
)

func init() {
	globalOpenIds = new(define.SecureWxOpenid)
	globalOpenIds.Data = make(map[string]*sgwxopenid.SWxOpenid)
}

func InitWxOpenIdCfg() {
	globalWxOpenIdCfg = new(sgwxdef.WxAppidCfg)
	cfgFile := sgcfg.GetServerCfgDir() + "wx_appid.json"
	err := sgcfg.ReadCfg(cfgFile, globalWxOpenIdCfg)
	if err != nil {
		sglog.Error("InitWxOpenIdCfg error,", err)
		sgthread.DelayExit(2)
	}

	globalAliAppidCfg = new(alisms.AliAppidCfg)
	cfgFile = sgcfg.GetServerCfgDir() + "ali_sms.json"
	err = sgcfg.ReadCfg(cfgFile, globalAliAppidCfg)
	if err != nil {
		sglog.Error("InitAliAppidIdCfg error,", err)
		sgthread.DelayExit(2)
	}

}

func AddWxOpenId(data *sgwxopenid.SWxOpenid) error {
	globalOpenIds.Lock.Lock()
	defer globalOpenIds.Lock.Unlock()

	if _, ok := globalOpenIds.Data[data.Code]; ok {
		return errors.New("code already exist")
	}
	globalOpenIds.Data[data.Code] = data
	return nil
}

func GetWxOpenId(code string) (*sgwxopenid.SWxOpenid, error) {
	globalOpenIds.Lock.Lock()
	defer globalOpenIds.Lock.Unlock()

	if v, ok := globalOpenIds.Data[code]; ok {
		return v, nil
	}
	return nil, errors.New("code not exisit")
}

func GetAppidCfg() (string, string) {
	return globalWxOpenIdCfg.Appid, globalWxOpenIdCfg.Secret
}
