package data

import (
	"xcxYaohao/src/config"

	"github.com/coderguang/GameEngine_go/sgcfg"

	"github.com/coderguang/GameEngine_go/sglog"
)

var (
	globalBoardcast *config.SecureBoardCast
)

func init() {
	globalBoardcast = new(config.SecureBoardCast)
	globalBoardcast.Data = make(map[string]string)
}

func ReloadBoardcast(cmd []string) {
	fileName := sgcfg.GetServerCfgDir() + "boardcast.json"
	cfgs := []config.BoardCastCfg{}
	err := sgcfg.ReadCfg(fileName, &cfgs)
	if err != nil {
		sglog.Error("ReloadBoardcast,err:", err)
		return
	}

	globalBoardcast.Lock.Lock()
	defer globalBoardcast.Lock.Unlock()
	for _, v := range cfgs {
		globalBoardcast.Data[v.Title] = v.Content
	}
}

func GetBoardcast(title string) string {
	globalBoardcast.Lock.Lock()
	defer globalBoardcast.Lock.Unlock()
	if v, ok := globalBoardcast.Data[title]; ok {
		return v
	}
	return "该城市暂不支持"
}
